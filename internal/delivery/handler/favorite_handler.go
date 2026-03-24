package handler

import (
	"context"
	"regexp"
	"stock_backend/internal/helper"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/request"
	"stock_backend/internal/service"
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
	Service   service.FavoriteService
	Validator *validator.Validate
}

func NewFavoriteHandler(service service.FavoriteService, validator *validator.Validate) FavoriteHandler {
	return &FavoriteHandlerImpl{
		Service:   service,
		Validator: validator,
	}
}

func (handler *FavoriteHandlerImpl) GetFavorites(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	res, err := handler.Service.GetFavorites(userId, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapFavoritesErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *FavoriteHandlerImpl) AddFavorites(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	var addFavoriteRequest request.AddFavoriteUnderwriterRequest
	if err := c.BodyParser(&addFavoriteRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrInvalidRequestBody.Error())
	}

	if err := handler.Validator.Struct(addFavoriteRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, helper.ValidationError(err))
	}

	res, err := handler.Service.CreateFavorite(userId, addFavoriteRequest.UnderwriterId, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapFavoritesErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (handler *FavoriteHandlerImpl) RemoveFavorites(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	underwriterCode := c.Params("underwriter")
	if !regexp.MustCompile(`^[A-Za-z]{2}$`).MatchString(underwriterCode) {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUnderwriterInvalid.Error())
	}

	res, err := handler.Service.RemoveFavorite(userId, underwriterCode, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapFavoritesErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
