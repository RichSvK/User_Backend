package router

import (
	"database/sql"
	"stock_backend/handler"
	"stock_backend/repository"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RegisterWatchlistRoutes(router fiber.Router, db *sql.DB) {
	validator := validator.New()
	watchlistRepository := repository.NewWatchlistRepository(db)
	watchlistService := service.NewWatchlistService(watchlistRepository)
	watchlistHandler := handler.NewWatchlistHandler(watchlistService, validator)

	authRouting := router.Group("/api/auth/watchlist")
	authRouting.Get("/", watchlistHandler.GetWatchlist)
	authRouting.Post("/", watchlistHandler.AddWatchlist)
	authRouting.Delete("/:stock", watchlistHandler.RemoveWatchlist)
}
