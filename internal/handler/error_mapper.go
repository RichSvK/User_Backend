package handler

import (
	domainerr "stock_backend/internal/model/domainerr"

	"github.com/gofiber/fiber/v2"
)

func MapErrorToHTTPStatus(err error) int {
	switch err {
	case domainerr.ErrUserNotFound:
		return fiber.StatusNotFound
	case domainerr.ErrWrongPassword:
		return fiber.StatusUnauthorized
	case domainerr.ErrNotVerified:
		return fiber.StatusForbidden
	case domainerr.ErrEmailExists:
		return fiber.StatusConflict
	case domainerr.ErrInvalidToken,
		domainerr.ErrInvalidTokenClaims,
		domainerr.ErrMissingSubject:
		return fiber.StatusUnauthorized
	default:
		return fiber.StatusInternalServerError
	}
}

func MapFavoritesErrorToHTTPStatus(err error) int {
	switch err {
	case domainerr.ErrFavoritesNotFound:
		return fiber.StatusNotFound
	case domainerr.ErrFavoritesDuplicate:
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}
}
