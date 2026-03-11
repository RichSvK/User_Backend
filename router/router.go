package router

import (
	"database/sql"
	"time"

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
		EnablePrintRoutes:     true,            // Print routes on startup
		BodyLimit:             4 * 1024 * 1024, // 4 MB request body limit
		Prefork:               true,           // Set to true to enable preforking
		CaseSensitive:         false,
		DisableStartupMessage: false,          // Disable Startup Message if needed
		JSONEncoder:           json.Marshal,   // Custom JSON Encoder
		JSONDecoder:           json.Unmarshal, // Custom JSON Decoder
		Views:                 nil,            // Set to nil if not using views
		StrictRouting:         true,
	})

	// Middleware Logger
	app.Use(logger.New())

	RegisterUserRoutes(app, db, redisDB)
	RegisterWatchlistRoutes(app, db)
	RegisterFavoriteRoutes(app, db, redisDB)

	return app
}
