package repository

import (
	"context"
	"database/sql"
	"errors"
	"stock_backend/model/entity"

	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetUser(email string, ctx context.Context) (*entity.User, error)
	Create(user entity.User, ctx context.Context) error
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
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) Create(user entity.User, ctx context.Context) error {
	query := "INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)"
	_, err := repository.DB.ExecContext(ctx, query, user.ID.String(), user.Username, user.Email, user.Password)
	return err
}

func (repository *UserRepositoryImpl) Logout(userId string, ctx context.Context) error {
	// Remove user favorites from Redis cache
	// return repository.RedisDB.Del(ctx, fmt.Sprintf("favorites:%s", userId)).Err()
	return nil
}

func (repository *UserRepositoryImpl) DeleteUser(userId string, ctx context.Context) error {
	query := "DELETE FROM users WHERE id = $1 AND roleid = 1"
	res, err := repository.DB.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	// return repository.RedisDB.Del(ctx, fmt.Sprintf("favorites:%s", userId)).Err()
	return nil
}

func (repository *UserRepositoryImpl) GetUserByID(userId string, ctx context.Context) (*entity.User, error) {
	query := "SELECT id, username, email FROM users WHERE id = $1"
	row := repository.DB.QueryRowContext(ctx, query, userId)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
