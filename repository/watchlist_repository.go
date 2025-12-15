package repository

import (
	"context"
	"database/sql"
	"log"
)

type WatchlistRepository interface {
	AddWatchlist(ctx context.Context, userId string, stock string) error
	RemoveWatchlist(ctx context.Context, userId string, stock string) error
	GetWatchlistByUserID(ctx context.Context, userId string) ([]string, error)
}

type WatchlistRepositoryImpl struct {
	DB *sql.DB
}

func NewWatchlistRepository(db *sql.DB) WatchlistRepository {
	return &WatchlistRepositoryImpl{
		DB: db,
	}
}

func (repository *WatchlistRepositoryImpl) AddWatchlist(ctx context.Context, userId string, stock string) error {
	query := "INSERT INTO watchlist (userid, stock) VALUES ($1, $2)"
	_, err := repository.DB.ExecContext(ctx, query, userId, stock)
	return err
}

func (repository *WatchlistRepositoryImpl) RemoveWatchlist(ctx context.Context, userId string, stock string) error {
	query := "DELETE FROM watchlist WHERE userid = $1 AND stock = $2"
	_, err := repository.DB.ExecContext(ctx, query, userId, stock)
	return err
}

func (repository *WatchlistRepositoryImpl) GetWatchlistByUserID(ctx context.Context, userId string) ([]string, error) {
	query := "SELECT stock FROM watchlist WHERE userid = $1"
	rows, err := repository.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var watchlist []string
	for rows.Next() {
		var stock string
		if err := rows.Scan(&stock); err != nil {
			return nil, err
		}
		watchlist = append(watchlist, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return watchlist, nil
}
