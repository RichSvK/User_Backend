package router

import (
	"stock_backend/handler"
	"stock_backend/repository"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterFavoriteRoutes(router fiber.Router, db *mongo.Client, redis_db *redis.Client) {
	favoriteRouting := router.Group("/api/auth/favorites")
	validator := validator.New()
	favoriteRepository := repository.NewFavoriteRepository(db, redis_db)
	favoriteService := service.NewFavoriteService(favoriteRepository)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService, validator)

	favoriteRouting.Get("/", favoriteHandler.GetFavorites)
	favoriteRouting.Post("/add", favoriteHandler.AddFavorites)
	favoriteRouting.Delete("/remove", favoriteHandler.RemoveFavorites)
}
