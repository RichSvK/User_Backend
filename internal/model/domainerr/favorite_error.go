package domainerr

import "errors"

var (
	ErrFavoritesNotFound           = errors.New("favorites not found")
	ErrFavoritesDuplicate          = errors.New("data already exists")
	ErrFavoritesUserIdRequired     = errors.New("user id is required")
	ErrFavoritesUnderwriterInvalid = errors.New("invalid underwriter code")
)
