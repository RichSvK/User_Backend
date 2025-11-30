package middleware

import (
	"stock_backend/model/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	// Check if the user is logged in
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Output{
			Message: "Internal server error",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if role != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(response.Output{
			Message: "Unauthorized access",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	return c.Next()
}

func UserMiddleware(c *fiber.Ctx) error {
	// Check if the user is logged in
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Output{
			Message: "Internal server error",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if role != "User" {
		return c.Status(fiber.StatusForbidden).JSON(response.Output{
			Message: "Unauthorized access",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	return c.Next()
}
