package handler

import (
	"context"
	"errors"
	"log"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/response"

	"github.com/gofiber/fiber/v2"
)

func ResponseErrorJSON(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(response.FailedResponse{
		Message: message,
	})
}

func MapErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, domainerr.ErrUserNotFound):
		return fiber.StatusNotFound, err.Error()

	case errors.Is(err, domainerr.ErrWrongPassword),
		errors.Is(err, domainerr.ErrInvalidToken),
		errors.Is(err, domainerr.ErrInvalidTokenClaims),
		errors.Is(err, domainerr.ErrMissingSubject):
		return fiber.StatusUnauthorized, err.Error()

	case errors.Is(err, domainerr.ErrNotVerified):
		return fiber.StatusForbidden, err.Error()

	case errors.Is(err, domainerr.ErrEmailExists),
		errors.Is(err, domainerr.ErrVerified):
		return fiber.StatusConflict, err.Error()

	case errors.Is(err, context.DeadlineExceeded):
		return fiber.StatusGatewayTimeout, err.Error()

	default:
		log.Printf("[ERROR] error: %v", err)
		return fiber.StatusInternalServerError, domainerr.ErrInternal.Error()
	}
}

func MapFavoritesErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, domainerr.ErrFavoritesNotFound):
		return fiber.StatusNotFound, err.Error()

	case errors.Is(err, domainerr.ErrFavoritesDuplicate):
		return fiber.StatusConflict, err.Error()

	case errors.Is(err, context.DeadlineExceeded):
		return fiber.StatusGatewayTimeout, err.Error()

	default:
		log.Printf("[ERROR] error: %v", err)
		return fiber.StatusInternalServerError, domainerr.ErrInternal.Error()
	}
}

func MapWatchlistErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, domainerr.ErrWatchlistNotFound):
		return fiber.StatusNotFound, err.Error()

	case errors.Is(err, domainerr.ErrWatchlistDuplicate):
		return fiber.StatusConflict, err.Error()

	case errors.Is(err, domainerr.ErrStockServiceUnavailable):
		return fiber.StatusServiceUnavailable, err.Error()

	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, domainerr.ErrServiceTimeout):
		return fiber.StatusGatewayTimeout, domainerr.ErrServiceTimeout.Error()

	default:
		log.Printf("[ERROR] error: %v", err)
		return fiber.StatusInternalServerError, domainerr.ErrInternal.Error()
	}
}
