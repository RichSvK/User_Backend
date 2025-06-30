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

func RegisterUserRoutes(router fiber.Router, db *sql.DB) {
	userRouting := router.Group("/users")
	validator := validator.New()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, validator)

	userRouting.Post("/login", userHandler.Login)
	userRouting.Post("/register", userHandler.Register)
	// userRouting.Post("/logout", middleware.JWTMiddleware, userHandler.Logout)
	// userRouting.Put("/update", middleware.JWTMiddleware, userHandler.UpdateUser)
}

func RegisterFavoriteRoutes(router fiber.Router, db *mongo.Client, redis_db *redis.Client) {
	favoriteRouting := router.Group("/favorites")
	validator := validator.New()
	favoriteRepository := repository.NewFavoriteRepository(db, redis_db)
	favoriteService := service.NewFavoriteService(favoriteRepository)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService, validator)

	favoriteRouting.Get("/", middleware.JWTMiddleware, favoriteHandler.GetFavorites)
	favoriteRouting.Post("/add", middleware.JWTMiddleware, favoriteHandler.AddFavorites)
}
