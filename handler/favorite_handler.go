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
	userId := c.Locals("userId").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	status, result := handler.Service.GetFavoritesService(userId, ctx)

	return c.Status(status).JSON(result)
}

func (handler *FavoriteHandlerImpl) AddFavorites(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var addFavoriteRequest request.FavoriteUnderwriterRequest
	if err := c.BodyParser(&addFavoriteRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	status, result := handler.Service.CreateFavorite(userId, addFavoriteRequest, ctx)

	return c.Status(status).JSON(result)
}

func (handler *FavoriteHandlerImpl) RemoveFavorites(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, result := handler.Service.RemoveFavoriteService(userId, ctx)

	return c.Status(status).JSON(result)
}
