package handler

import (
	"context"
	"stock_backend/model/request"
	"stock_backend/model/response"
	"stock_backend/service"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FavoriteHandler interface {
	GetFavorites(c *fiber.Ctx) error
	AddFavorites(c *fiber.Ctx) error
	RemoveFavorites(c *fiber.Ctx) error
}

type FavoriteHandlerImpl struct {
	Service service.FavoriteService
}

func NewFavoriteHandler(service service.FavoriteService, validator *validator.Validate) FavoriteHandler {
	return &FavoriteHandlerImpl{
		Service: service,
	}
}

func (handler *FavoriteHandlerImpl) GetFavorites(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "User ID is required",
		})
	}

	res, err := handler.Service.GetFavorites(userId, ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *FavoriteHandlerImpl) AddFavorites(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "User ID is required",
		})
	}

	var addFavoriteRequest request.FavoriteUnderwriterRequest
	if err := c.BodyParser(&addFavoriteRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Invalid request",
		})
	}

	res, err := handler.Service.CreateFavorite(userId, addFavoriteRequest.UnderwriterId, ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (handler *FavoriteHandlerImpl) RemoveFavorites(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "User ID is required",
		})
	}

	underwriterCode := c.Params("underwriter")
	if underwriterCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Underwriter code is required",
		})
	}

	res, err := handler.Service.RemoveFavorite(userId, underwriterCode, ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
