package domain_error

import "errors"

var (
	ErrWatchlistNotFound = errors.New("watchlist not found")
)
