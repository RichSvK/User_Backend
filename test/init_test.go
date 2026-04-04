package test

import (
	"database/sql"
	"os"
	"stock_backend/config"
	"stock_backend/internal/delivery/router"
	"testing"

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

func TestMain(m *testing.M) {
	ClearTable("users")
	if err := CreateTestUser(email, password); err != nil {
		panic(err)
	}

	tok, err := GetUserToken(email, password)
	if err != nil {
		panic(err)
	}
	token = tok

	if err := CreateAdmin("admin@gmail.com", password, "admin"); err != nil {
		panic(err)
	}

	adminToken, err = GetUserToken("admin@gmail.com", password)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}
