package router

import (
	"database/sql"
	"log"
	"os"
	"stock_backend/internal/delivery/handler"
	"stock_backend/internal/delivery/middleware"
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

	smtp, err := service.LoadSMTPConfig()
	if err != nil {
		log.Println("failed to load smtp")
	}
	userService := service.NewUserService(userRepository, jwtSecret, smtp)
	userHandler := handler.NewUserHandler(userService, validator)

	userRouting := router.Group("/api/v1/users")
	userRouting.Use(middleware.LoggedOutMiddleware())
	userRouting.Post("/login", userHandler.Login)
	userRouting.Post("/register", userHandler.Register)
	userRouting.Get("/verify", userHandler.VerifyUser)

	authRouting := router.Group("/api/v1/auth/users")
	authRouting.Use(middleware.JWTMiddleware(jwtSecret))
	authRouting.Get("/profile", userHandler.GetUserInfo)
	authRouting.Post("/logout", userHandler.Logout)

	adminRouting := authRouting.Use(middleware.AdminMiddleware())
	adminRouting.Delete("", userHandler.DeleteUser)
}
