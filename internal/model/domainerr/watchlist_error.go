package domainerr

import "errors"

var (
	ErrWatchlistNotFound  = errors.New("watchlist not found")
	ErrWatchlistDuplicate = errors.New("duplicate stock in watchlist")
)