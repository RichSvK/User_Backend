package handler

import (
	"context"
	"errors"
	"stock_backend/internal/helper"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"stock_backend/internal/service"
	"time"

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
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	res, err := handler.Service.GetWatchlist(ctx, userId)
	if err != nil {
		status, message := MapWatchlistErrorToHTTPStatus(err)
		return ResponseErrorJSON(c, status, message)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *WatchlistHandlerImpl) AddWatchlist(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	var req request.AddWatchlistRequest
	if err := c.BodyParser(&req); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrInvalidRequestBody.Error())
	}

	if err := handler.Validator.Struct(req); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, helper.ValidationError(err))
	}

	res, err := handler.Service.AddToWatchlist(ctx, userId, req.Stock)
	if err != nil {
		var serviceErr *domainerr.ServiceError
		if errors.As(err, &serviceErr) {
			return c.Status(serviceErr.Code).JSON(response.FailedResponse{
				Message: serviceErr.Message,
			})
		}
		status, message := MapWatchlistErrorToHTTPStatus(err)

		return ResponseErrorJSON(c, status, message)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (handler *WatchlistHandlerImpl) RemoveWatchlist(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	stock := c.Params("stock")
	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	res, err := handler.Service.RemoveFromWatchlist(ctx, userId, stock)
	if err != nil {
		status, message := MapWatchlistErrorToHTTPStatus(err)
		return ResponseErrorJSON(c, status, message)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
