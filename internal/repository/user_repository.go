package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"stock_backend/internal/entity"
	"stock_backend/internal/model/domainerr"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetUser(email string, ctx context.Context) (*entity.User, error)
	Create(user entity.User, ctx context.Context) (*entity.User, error)
	VerifyUser(userId string, ctx context.Context) error
	Logout(userId string, ctx context.Context) error
	DeleteUser(userId string, ctx context.Context) error
	GetUserByID(userId string, ctx context.Context) (*entity.User, error)
}

type UserRepositoryImpl struct {
	DB      *sql.DB
	RedisDB *redis.Client
}

func NewUserRepository(db *sql.DB, redis_db *redis.Client) UserRepository {
	return &UserRepositoryImpl{
		DB:      db,
		RedisDB: redis_db,
	}
}

func (repository *UserRepositoryImpl) GetUser(email string, ctx context.Context) (*entity.User, error) {
	query := "SELECT id, username, email, password, r.rolename, verified FROM users u JOIN roles r ON u.roleid = r.roleid WHERE email = $1"
	row := repository.DB.QueryRowContext(ctx, query, email)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Verified)
	if err == sql.ErrNoRows {
		return nil, domainerr.ErrUserNotFound
	}

	if err != nil {
		return nil, domainerr.ErrInternal
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) Create(user entity.User, ctx context.Context) (*entity.User, error) {
	tx, err := repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Println("transaction rollback error:", err)
		}
	}()

	insertQuery := `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, verified;
	`
	var createdUser entity.User
	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		user.ID.String(),
		user.Username,
		user.Email,
		user.Password,
	).Scan(&createdUser.ID, &createdUser.Username, &createdUser.Email, &createdUser.Verified)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, domainerr.ErrEmailExists
			}
		}
		return nil, domainerr.ErrInternal
	}

	return &createdUser, tx.Commit()
}

func (repository *UserRepositoryImpl) VerifyUser(userId string, ctx context.Context) error {
	query := "UPDATE users SET verified = TRUE WHERE id = $1 AND verified = FALSE"
	res, err := repository.DB.ExecContext(ctx, query, userId)
	if err != nil {
		return domainerr.ErrInternal
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return domainerr.ErrInternal
	}

	fmt.Printf("Rows affected: %d\n", rows) // Add this
	if rows == 0 {
		return domainerr.ErrVerified
	}
	return nil
}

func (repository *UserRepositoryImpl) Logout(userId string, ctx context.Context) error {
	// Remove user favorites from Redis cache
	if repository.RedisDB != nil {
		if err := repository.RedisDB.Del(
			ctx,
			fmt.Sprintf("favorites:%s", userId),
		).Err(); err != nil {
			return domainerr.ErrInternal
		}
	}

	return nil
}

func (repository *UserRepositoryImpl) DeleteUser(userId string, ctx context.Context) error {
	query := "DELETE FROM users WHERE id = $1"
	res, err := repository.DB.ExecContext(ctx, query, userId)
	if err != nil {
		return domainerr.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return domainerr.ErrInternal
	}

	if rowsAffected == 0 {
		return domainerr.ErrUserNotFound
	}

	// Remove user favorites from Redis cache
	if repository.RedisDB != nil {
		_ = repository.RedisDB.Del(ctx, fmt.Sprintf("favorites:%s", userId)).Err()
	}
	return nil
}

func (repository *UserRepositoryImpl) GetUserByID(userId string, ctx context.Context) (*entity.User, error) {
	query := "SELECT id, username, email FROM users WHERE id = $1"
	row := repository.DB.QueryRowContext(ctx, query, userId)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		return nil, domainerr.ErrUserNotFound
	}

	if err != nil {
		return nil, domainerr.ErrInternal
	}

	return &user, nil
}
