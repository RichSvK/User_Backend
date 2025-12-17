package handler

import (
	"stock_backend/model/request"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type WatchlistHandler interface {
	GetWatchlist(c *fiber.Ctx) error
	AddWatchlist(c *fiber.Ctx) error
	RemoveWatchlist(c *fiber.Ctx) error
}

type WatchlistHandlerImpl struct {
	Service   service.WatchlistService
	Validator *validator.Validate
}

func NewWatchlistHandler(service service.WatchlistService, validator *validator.Validate) WatchlistHandler {
	return &WatchlistHandlerImpl{
		Service:   service,
		Validator: validator,
	}
}

func (handler *WatchlistHandlerImpl) GetWatchlist(c *fiber.Ctx) error {
	userId := c.Get("X-User-ID")
	status, result := handler.Service.GetWatchlist(c.Context(), userId)
	return c.Status(status).JSON(result)
}

func (handler *WatchlistHandlerImpl) AddWatchlist(c *fiber.Ctx) error {
	userId := c.Get("X-User-ID")

	var req request.AddWatchlistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Stock == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Stock is required",
		})
	}

	status, result := handler.Service.AddToWatchlist(c.Context(), userId, req.Stock)
	return c.Status(status).JSON(result)
}

func (handler *WatchlistHandlerImpl) RemoveWatchlist(c *fiber.Ctx) error {
	userId := c.Get("X-User-ID")
	stock := c.Params("stock")
	status, result := handler.Service.RemoveFromWatchlist(c.Context(), userId, stock)
	return c.Status(status).JSON(result)
}
