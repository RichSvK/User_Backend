package middleware

import (
	"stock_backend/internal/handler"
	"stock_backend/internal/model/domainerr"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the user is logged in
		role, ok := c.Locals("role").(string)
		if !ok {
			return handler.ResponseErrorJSON(c, fiber.StatusInternalServerError, domainerr.ErrInternal.Error())
		}

		if role != "admin" {
			return handler.ResponseErrorJSON(c, fiber.StatusForbidden, domainerr.ErrUnauthorizedAccess.Error())
		}
		return c.Next()
	}
}

func UserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the user is logged in
		role, ok := c.Locals("role").(string)
		if !ok {
			return handler.ResponseErrorJSON(c, fiber.StatusInternalServerError, domainerr.ErrInternal.Error())
		}

		if role != "user" {
			return handler.ResponseErrorJSON(c, fiber.StatusForbidden, domainerr.ErrUnauthorizedAccess.Error())
		}

		return c.Next()
	}
}
