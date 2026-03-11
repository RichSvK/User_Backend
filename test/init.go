package test

import (
	"database/sql"
	"stock_backend/database"
	"stock_backend/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var app *fiber.App
var db *sql.DB

func init() {
	_ = godotenv.Load("../test.env")

	db = database.DatabaseConfig()
	redisDb := database.ConnectRedis()
	app = router.SetupRouter(db, redisDb)
}
