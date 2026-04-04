package handler

import (
	"context"
	"time"

	"stock_backend/internal/helper"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/request"
	"stock_backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	VerifyUser(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
}

type UserHandlerImpl struct {
	UserService service.UserService
	Validator   *validator.Validate
}

func NewUserHandler(service service.UserService, validator *validator.Validate) UserHandler {
	return &UserHandlerImpl{
		UserService: service,
		Validator:   validator,
	}
}

func (handler *UserHandlerImpl) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrInvalidRequestBody.Error())
	}

	if err := handler.Validator.Struct(loginRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, helper.ValidationError(err))
	}

	res, err := handler.UserService.Login(loginRequest, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *UserHandlerImpl) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var registerRequest request.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrInvalidRequestBody.Error())
	}

	if err := handler.Validator.Struct(registerRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, helper.ValidationError(err))
	}

	res, err := handler.UserService.Register(registerRequest, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (handler *UserHandlerImpl) VerifyUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	token := c.Query("token")
	if token == "" {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrEmptyToken.Error())
	}

	res, err := handler.UserService.VerifyUser(token, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *UserHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	res, err := handler.UserService.Logout(userId, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *UserHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var deleteRequest request.DeleteUserRequest
	if err := c.BodyParser(&deleteRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrInvalidRequestBody.Error())
	}

	if err := handler.Validator.Struct(deleteRequest); err != nil {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, helper.ValidationError(err))
	}

	res, err := handler.UserService.DeleteUser(deleteRequest.UserId, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (handler *UserHandlerImpl) GetUserInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId, ok := helper.GetUserID(c)
	if !ok {
		return ResponseErrorJSON(c, fiber.StatusBadRequest, domainerr.ErrFavoritesUserIdRequired.Error())
	}

	res, err := handler.UserService.GetProfile(userId, ctx)
	if err != nil {
		return ResponseErrorJSON(c, MapErrorToHTTPStatus(err), err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
