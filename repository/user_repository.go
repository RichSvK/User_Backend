package repository

import (
	"context"
	"database/sql"
	"stock_backend/model/entity"
)

type UserRepository interface {
	GetUser(email string, ctx context.Context) (*entity.User, error)
	Create(user entity.User, ctx context.Context) error
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (repository *UserRepositoryImpl) GetUser(email string, ctx context.Context) (*entity.User, error) {
	query := "SELECT username, email, password FROM users WHERE email = $1"
	row := repository.DB.QueryRowContext(ctx, query, email)

	var user entity.User
	err := row.Scan(&user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserRepositoryImpl) Create(user entity.User, ctx context.Context) error {
	query := "INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)"
	_, err := repository.DB.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password)
	return err
}
