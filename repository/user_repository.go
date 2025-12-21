package repository

import (
	"context"
	"database/sql"
	"fmt"
	"stock_backend/model/entity"
	domain_error "stock_backend/model/error"

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
	query := "SELECT id, username, email, password, r.rolename FROM users u JOIN roles r ON u.roleid = r.roleid WHERE email = $1"
	row := repository.DB.QueryRowContext(ctx, query, email)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return nil, domain_error.ErrUserNotFound
	}

	if err != nil {
		return nil, domain_error.ErrInternal
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) Create(user entity.User, ctx context.Context) (*entity.User, error) {
	tx, err := repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var resultUser entity.User

	// Try to update if user exists and NOT verified
	updateQuery := `
		UPDATE users
		SET username = $1,
		    password = $2
		WHERE email = $3 AND verified = false
		RETURNING id, username, email;
	`

	err = tx.QueryRowContext(
		ctx,
		updateQuery,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&resultUser.ID,
		&resultUser.Username,
		&resultUser.Email,
	)

	// Update succeeded, row returned
	if err == nil {
		err = tx.Commit()
		return &resultUser, err
	}

	// If error is not "no rows", return it
	if err != sql.ErrNoRows {
		return nil, domain_error.ErrInternal
	}

	// Otherwise insert new user
	insertQuery := `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email;
	`

	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		user.ID.String(),
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&resultUser.ID,
		&resultUser.Username,
		&resultUser.Email,
	)

	if err != nil {
		return nil, domain_error.ErrInternal
	}

	return &resultUser, tx.Commit()
}

func (repository *UserRepositoryImpl) VerifyUser(userId string, ctx context.Context) error {
	query := "UPDATE users SET verified = TRUE WHERE id = $1"
	res, err := repository.DB.ExecContext(ctx, query, userId)
	if err != nil {
		return domain_error.ErrInternal
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain_error.ErrUserNotFound
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
			return domain_error.ErrInternal
		}
	}

	return nil
}

func (repository *UserRepositoryImpl) DeleteUser(userId string, ctx context.Context) error {
	query := "DELETE FROM users WHERE id = $1"
	res, err := repository.DB.ExecContext(ctx, query, userId)
	if err != nil {
		return domain_error.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return domain_error.ErrInternal
	}

	if rowsAffected == 0 {
		return domain_error.ErrUserNotFound
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
		return nil, domain_error.ErrUserNotFound
	}

	if err != nil {
		return nil, domain_error.ErrInternal
	}

	return &user, nil
}
