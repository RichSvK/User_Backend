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
)

func RegisterWatchlistRoutes(router fiber.Router, db *sql.DB) {
	validator := validator.New()
	watchlistRepository := repository.NewWatchlistRepository(db)
	watchlistService := service.NewWatchlistService(watchlistRepository)
	watchlistHandler := handler.NewWatchlistHandler(watchlistService, validator)

	authRouting := router.Group("/api/v1/auth/watchlist")
	authRouting.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
	authRouting.Get("", watchlistHandler.GetWatchlist)
	authRouting.Post("", watchlistHandler.AddWatchlist)
	authRouting.Delete("/:stock", watchlistHandler.RemoveWatchlist)
}
