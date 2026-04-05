package router

import (
	"database/sql"
	"stock_backend/internal/delivery/middleware"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(db *sql.DB, redisDB *redis.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Stock Backend API",
		IdleTimeout:           5 * time.Second,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		EnablePrintRoutes:     true,
		BodyLimit:             4 * 1024 * 1024,
		Prefork:               false, // Set to true to enable preforking
		CaseSensitive:         false,
		DisableStartupMessage: false,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		Views:                 nil,
		StrictRouting:         true,
	})

	// Middleware setup
	app.Use(logger.New())
	middleware.CorsMiddleware(app)
	middleware.RateLimitMiddleware(app)

	validator := validator.New()

	// Register Route
	RegisterUserRoutes(app, db, validator, redisDB)
	RegisterWatchlistRoutes(app, db, validator)
	RegisterFavoriteRoutes(app, db, validator, redisDB)
	return app
}
