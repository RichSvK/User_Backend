package repository

import (
	"context"
	"database/sql"
	"log"
	"stock_backend/internal/model/domainerr"
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

	if err != nil {
		return domainerr.ErrInternal
	}

	return nil
}

func (repository *WatchlistRepositoryImpl) RemoveWatchlist(ctx context.Context, userId string, stock string) error {
	query := "DELETE FROM watchlist WHERE userid = $1 AND stock = $2"
	res, err := repository.DB.ExecContext(ctx, query, userId, stock)
	if err != nil {
		return domainerr.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return domainerr.ErrInternal
	}

	if rowsAffected == 0 {
		return domainerr.ErrWatchlistNotFound
	}

	return nil
}

func (repository *WatchlistRepositoryImpl) GetWatchlistByUserID(ctx context.Context, userId string) ([]string, error) {
	query := "SELECT stock FROM watchlist WHERE userid = $1"
	rows, err := repository.DB.QueryContext(ctx, query, userId)

	if err != nil {
		log.Printf("query watchlist failed: %v\n", err)
		return nil, domainerr.ErrInternal
	}

	defer func() {
		_ = rows.Close()
	}()

	var watchlist []string
	for rows.Next() {
		var stock string
		if err := rows.Scan(&stock); err != nil {
			log.Printf("scan watchlist row failed: %v", err)
			return nil, domainerr.ErrInternal
		}
		watchlist = append(watchlist, stock)
	}

	if err := rows.Err(); err != nil {
		log.Printf("row iteration failed: %v", err)
		return nil, domainerr.ErrInternal
	}

	return watchlist, nil
}
