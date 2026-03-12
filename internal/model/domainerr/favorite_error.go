package domainerr

import "errors"

var (
	ErrFavoritesNotFound  = errors.New("favorites not found")
	ErrFavoritesDuplicate = errors.New("data already exists")
)
