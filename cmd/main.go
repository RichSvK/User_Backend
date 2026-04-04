package main

import (
	"log"
	"os"
	"stock_backend/config"
	"stock_backend/internal/delivery/router"
)

func main() {
	// Load local environment variables => not in production
	config.LoadEnv(".env")

	// Connect to database
	db := config.DatabaseConfig()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close DB: %v\n", err)
		}
	}()

	redisDb := config.ConnectRedis()
	defer func() {
		if err := redisDb.Close(); err != nil {
			log.Printf("Failed to close Redis client DB: %v\n", err)
		}
	}()

	// Routes Grouping
	app := router.SetupRouter(db, redisDb)

	// Run the app
	if err := app.Listen(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err.Error())
	}
}
