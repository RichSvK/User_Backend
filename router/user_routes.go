package router

import (
	"database/sql"
	"os"
	"stock_backend/internal/handler"
	"stock_backend/internal/middleware"
	"stock_backend/internal/repository"
	"stock_backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRoutes(router fiber.Router, db *sql.DB, redis_db *redis.Client) {
	validator := validator.New()
	jwtSecret := os.Getenv("JWT_SECRET")
	userRepository := repository.NewUserRepository(db, redis_db)
	userService := service.NewUserService(userRepository, os.Getenv("EMAIL_SECRET_KEY"), jwtSecret)
	userHandler := handler.NewUserHandler(userService, validator)

	userRouting := router.Group("/api/v1/users")
	userRouting.Post("/login", userHandler.Login)
	userRouting.Post("/register", userHandler.Register)
	userRouting.Get("/verify", userHandler.VerifyUser)

	authRouting := router.Group("/api/v1/auth/users")
	authRouting.Use(middleware.JWTMiddleware(jwtSecret))
	authRouting.Get("/profile", userHandler.GetUserInfo)
	authRouting.Post("/logout", userHandler.Logout)
	authRouting.Delete("/delete", userHandler.DeleteUser)
}
