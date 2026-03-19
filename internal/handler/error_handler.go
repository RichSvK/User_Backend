package handler

import (
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/response"

	"github.com/gofiber/fiber/v2"
)

func ResponseErrorJSON(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(response.FailedResponse{
		Message: message,
	})
}

func MapErrorToHTTPStatus(err error) int {
	switch err {
	case domainerr.ErrUserNotFound:
		return fiber.StatusNotFound
	case domainerr.ErrWrongPassword:
		return fiber.StatusUnauthorized
	case domainerr.ErrNotVerified:
		return fiber.StatusForbidden
	case domainerr.ErrEmailExists, domainerr.ErrVerified:
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

func MapWatchlistErrorToHTTPStatus(err error) int {
	switch err {
	case domainerr.ErrWatchlistNotFound:
		return fiber.StatusNotFound
	case domainerr.ErrWatchlistDuplicate:
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}
}
