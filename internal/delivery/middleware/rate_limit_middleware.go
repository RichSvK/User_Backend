package middleware

import (
	"stock_backend/internal/delivery/handler"
	"stock_backend/internal/model/domainerr"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimitMiddleware(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 60 * time.Second,

		LimitReached: func(c *fiber.Ctx) error {
			return handler.ResponseErrorJSON(c, fiber.StatusTooManyRequests, domainerr.ErrTooManyRequest.Error())
		},
	}))
}
