package middleware

import (
	"log"
	"os"
	"stock_backend/internal/delivery/handler"
	"stock_backend/internal/model/domainerr"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/gofiber/storage/redis/v3"
)

func RateLimitMiddleware(app *fiber.App) {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal("dailed to load get port")
	}

	store := redis.New(redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Port:     port,
	})

	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 60 * time.Second,
		Storage:    store,
		LimitReached: func(c *fiber.Ctx) error {
			return handler.ResponseErrorJSON(c, fiber.StatusTooManyRequests, domainerr.ErrTooManyRequest.Error())
		},
	}))
}
