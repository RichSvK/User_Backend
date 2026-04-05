package test

import (
	"database/sql"
	"stock_backend/config"
	"stock_backend/internal/delivery/router"

	"github.com/gofiber/fiber/v2"
)

var app *fiber.App
var db *sql.DB

const email = "richardsugiharto0@gmail.com"
const password = "87654321"

var token string
var adminToken string

func init() {
	config.LoadEnv("../test.env")

	db = config.DatabaseConfig()
	redisDb := config.ConnectRedis()
	app = router.SetupRouter(db, redisDb)
}
