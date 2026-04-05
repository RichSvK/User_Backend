package router

import (
	"database/sql"
	"os"
	"stock_backend/internal/circuit"
	"stock_backend/internal/client"
	"stock_backend/internal/delivery/handler"
	"stock_backend/internal/delivery/middleware"
	"stock_backend/internal/repository"
	"stock_backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RegisterWatchlistRoutes(router fiber.Router, db *sql.DB, validator *validator.Validate) {
	watchlistRepository := repository.NewWatchlistRepository(db)
	breaker := circuit.NewCircuitBreaker("stock-service")
	stockClient := client.NewStockClient(os.Getenv("STOCK_SERVICE_URL"), breaker)
	watchlistService := service.NewWatchlistService(watchlistRepository, stockClient)
	watchlistHandler := handler.NewWatchlistHandler(watchlistService, validator)

	authRouting := router.Group("/api/v1/watchlists")
	authRouting.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
	authRouting.Get("", watchlistHandler.GetWatchlist)
	authRouting.Post("/stocks", watchlistHandler.AddWatchlist)
	authRouting.Delete("/stocks/:stock", watchlistHandler.RemoveWatchlist)
}
