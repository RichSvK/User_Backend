package middleware

import (
	"stock_backend/model/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimitMiddleware(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 60 * time.Second,

		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(response.FailedResponse{
				Message: "Too many requests",
			})
		},
	}))
}
