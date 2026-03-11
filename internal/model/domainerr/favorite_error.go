package domainerr

import "errors"

var (
	ErrFavoritesNotFound  = errors.New("Favorites not found")
	ErrFavoritesDuplicate = errors.New("Data already exists")
)
