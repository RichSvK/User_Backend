package main

import (
	"log"
	"os"
	"stock_backend/config"
	"stock_backend/database"
	"stock_backend/router"
)

func main() {
	// Load local environment variables => not in production
	config.LoadEnv()

	// Connect to database
	db := database.DatabaseConfig()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close DB: %v\n", err)
		}
	}()

	redisDb := database.ConnectRedis()
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
