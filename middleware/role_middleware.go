package middleware

import (
	"stock_backend/model/response"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	// Check if the user is logged in
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
			Message: "Internal server error",
		})
	}

	if role != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(response.FailedResponse{
			Message: "Unauthorized access",
		})
	}

	return c.Next()
}

func UserMiddleware(c *fiber.Ctx) error {
	// Check if the user is logged in
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
			Message: "Internal server error",
		})
	}

	if role != "User" {
		return c.Status(fiber.StatusForbidden).JSON(response.FailedResponse{
			Message: "Unauthorized access",
		})
	}

	return c.Next()
}
