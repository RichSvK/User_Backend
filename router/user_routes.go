package router

import (
	"database/sql"
	"stock_backend/handler"
	"stock_backend/repository"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRoutes(router fiber.Router, db *sql.DB, redis_db *redis.Client) {
	validator := validator.New()
	userRepository := repository.NewUserRepository(db, redis_db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, validator)

	userRouting := router.Group("/api/user")
	userRouting.Post("/login", userHandler.Login)
	userRouting.Post("/register", userHandler.Register)
	userRouting.Get("/verify", userHandler.VerifyUser)

	authRouting := router.Group("/api/auth/user")
	authRouting.Get("/profile", userHandler.GetUserInfo)
	authRouting.Post("/logout", userHandler.Logout)
	authRouting.Delete("/delete", userHandler.DeleteUser)
}
