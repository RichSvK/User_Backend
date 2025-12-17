package main

import (
	"log"
	"stock_backend/database"
	"stock_backend/router"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:               "Stock Backend API",
		IdleTimeout:           5 * time.Second,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		EnablePrintRoutes:     true,            // Print routes on startup
		BodyLimit:             4 * 1024 * 1024, // 4 MB request body limit
		Prefork:               false,           // Set to true for production to enable preforking
		CaseSensitive:         false,
		DisableStartupMessage: false,          // Disable Startup Message if needed
		JSONEncoder:           json.Marshal,   // Custom JSON Encoder
		JSONDecoder:           json.Unmarshal, // Custom JSON Decoder
		Views:                 nil,            // Set to nil if not using views
	})

	// Middleware
	app.Use(logger.New())

	// Load local environment variables => not in production
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading local .env file")
	}

	db, err := database.DatabaseConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close DB: %v\n", err)
		}
	}()

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// db_favorite, err := database.ConnectMongoDB(ctx)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// redisDb, err := database.ConnectRedis()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// Routes Grouping
	router.RegisterUserRoutes(app, db, nil)
	router.RegisterWatchlistRoutes(app, db)
	// router.RegisterFavoriteRoutes(app, db_favorite, redisDb)

	// Run the app
	if err := app.Listen(":8888"); err != nil {
		log.Fatal(err.Error())
	}
}
