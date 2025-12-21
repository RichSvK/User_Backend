package handler

import (
	domain_error "stock_backend/model/error"

	"github.com/gofiber/fiber/v2"
)

func MapErrorToHTTPStatus(err error) int {
	switch err {
	case domain_error.ErrUserNotFound:
		return fiber.StatusNotFound
	case domain_error.ErrWrongPassword:
		return fiber.StatusUnauthorized
	case domain_error.ErrNotVerified:
		return fiber.StatusForbidden
	case domain_error.ErrEmailExists:
		return fiber.StatusConflict
	case domain_error.ErrInvalidToken,
		domain_error.ErrInvalidTokenClaims,
		domain_error.ErrMissingSubject:
		return fiber.StatusUnauthorized
	default:
		return fiber.StatusInternalServerError
	}
}
