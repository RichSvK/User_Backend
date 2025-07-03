package router

import (
	"database/sql"
	"stock_backend/handler"
	"stock_backend/middleware"
	"stock_backend/repository"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterUserRoutes(router fiber.Router, db *sql.DB, redis_db *redis.Client) {
	userRouting := router.Group("/users")
	validator := validator.New()
	userRepository := repository.NewUserRepository(db, redis_db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, validator)

	userRouting.Post("/login", middleware.LoggedOutMiddleware, userHandler.Login)
	userRouting.Post("/register", middleware.LoggedOutMiddleware, userHandler.Register)
	userRouting.Post("/logout", middleware.LoggedInStatusMiddleware, middleware.JWTMiddleware, userHandler.Logout)
}

func RegisterFavoriteRoutes(router fiber.Router, db *mongo.Client, redis_db *redis.Client) {
	favoriteRouting := router.Group("/favorites")
	validator := validator.New()
	favoriteRepository := repository.NewFavoriteRepository(db, redis_db)
	favoriteService := service.NewFavoriteService(favoriteRepository)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService, validator)

	favoriteRouting.Get("/", middleware.LoggedInStatusMiddleware, middleware.JWTMiddleware, favoriteHandler.GetFavorites)
	favoriteRouting.Post("/add", middleware.LoggedInStatusMiddleware, middleware.JWTMiddleware, favoriteHandler.AddFavorites)
}
