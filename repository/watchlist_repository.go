package repository

import (
	"context"
	"database/sql"
	"log"
	domain_error "stock_backend/model/error"
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
		return domain_error.ErrInternal
	}

	return nil
}

func (repository *WatchlistRepositoryImpl) RemoveWatchlist(ctx context.Context, userId string, stock string) error {
	query := "DELETE FROM watchlist WHERE userid = $1 AND stock = $2"
	res, err := repository.DB.ExecContext(ctx, query, userId, stock)
	if err != nil {
		return domain_error.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return domain_error.ErrInternal
	}

	if rowsAffected == 0 {
		return domain_error.ErrWatchlistNotFound
	}

	return nil
}

func (repository *WatchlistRepositoryImpl) GetWatchlistByUserID(ctx context.Context, userId string) ([]string, error) {
	query := "SELECT stock FROM watchlist WHERE userid = $1"
	rows, err := repository.DB.QueryContext(ctx, query, userId)

	if err != nil {
		log.Printf("query watchlist failed: %v\n", err)
		return nil, domain_error.ErrInternal
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("failed to close rows: %v\n", cerr)
		}
	}()

	var watchlist []string
	for rows.Next() {
		var stock string
		if err := rows.Scan(&stock); err != nil {
			log.Printf("scan watchlist row failed: %v", err)
			return nil, domain_error.ErrInternal
		}
		watchlist = append(watchlist, stock)
	}

	if err := rows.Err(); err != nil {
		log.Printf("row iteration failed: %v", err)
		return nil, domain_error.ErrInternal
	}

	return watchlist, nil
}
