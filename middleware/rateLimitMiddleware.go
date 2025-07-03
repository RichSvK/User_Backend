package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimitMiddleware(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 60 * time.Second,

		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too much request!",
			})
		},
	}))
}
