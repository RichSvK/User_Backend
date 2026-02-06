package handler

import (
	"stock_backend/model/request"
	"stock_backend/model/response"
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
	res, err := handler.Service.GetWatchlist(c.Context(), userId)
	if err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *WatchlistHandlerImpl) AddWatchlist(c *fiber.Ctx) error {
	userId := c.Get("X-User-ID")

	var req request.AddWatchlistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Invalid request",
		})
	}

	if req.Stock == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Stock is required",
		})
	}

	res, err := handler.Service.AddToWatchlist(c.Context(), userId, req.Stock)
	if err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (handler *WatchlistHandlerImpl) RemoveWatchlist(c *fiber.Ctx) error {
	userId := c.Get("X-User-ID")
	stock := c.Params("stock")

	res, err := handler.Service.RemoveFromWatchlist(c.Context(), userId, stock)
	if err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
