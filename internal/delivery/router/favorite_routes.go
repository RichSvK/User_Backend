package router

import (
	"database/sql"
	"os"
	"stock_backend/internal/delivery/handler"
	"stock_backend/internal/delivery/middleware"
	"stock_backend/internal/repository"
	"stock_backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterFavoriteRoutes(router fiber.Router, db *sql.DB, validator *validator.Validate, redis_db *redis.Client) {
	favoriteRepository := repository.NewFavoriteRepository(db, redis_db)
	favoriteService := service.NewFavoriteService(favoriteRepository)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService, validator)

	favoriteRouting := router.Group("/api/v1/favorites")
	favoriteRouting.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")), middleware.UserMiddleware())
	favoriteRouting.Get("", favoriteHandler.GetFavorites)
	favoriteRouting.Post("", favoriteHandler.AddFavorites)
	favoriteRouting.Delete("/:underwriter", favoriteHandler.RemoveFavorites)
}
