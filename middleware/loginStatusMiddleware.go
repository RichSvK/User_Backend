package middleware

import (
	"stock_backend/model/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggedInStatusMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "You are logged out",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	return c.Next()
}

func LoggedOutMiddleware(c *fiber.Ctx) error {
	if c.Cookies("token") != "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Already logged in",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	return c.Next()
}
