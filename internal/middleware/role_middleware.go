package middleware

import (
	"stock_backend/internal/model/response"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the user is logged in
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
				Message: "Internal server error",
			})
		}

		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(response.FailedResponse{
				Message: "Unauthorized access",
			})
		}
		return c.Next()
	}
}

func UserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the user is logged in
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
				Message: "Internal server error",
			})
		}

		if role != "user" {
			return c.Status(fiber.StatusForbidden).JSON(response.FailedResponse{
				Message: "Unauthorized access",
			})
		}

		return c.Next()
	}
}
