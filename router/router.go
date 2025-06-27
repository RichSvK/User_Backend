package router

import (
	"database/sql"
	"stock_backend/handler"
	"stock_backend/repository"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(router fiber.Router, db *sql.DB) {
	userRouting := router.Group("/users")
	validator := validator.New()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, validator)

	userRouting.Post("/login", userHandler.Login)
	userRouting.Post("/register", userHandler.Register)
}
