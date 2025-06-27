package handler

import (
	"context"
	"time"

	"stock_backend/model/request"
	"stock_backend/model/response"
	"stock_backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Login(ctx *fiber.Ctx) error
	Register(c *fiber.Ctx) error
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

func (userHandler *UserHandlerImpl) Login(c *fiber.Ctx) error {
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if err := userHandler.Validator.Struct(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Validation failed",
			Time:    time.Now(),
			Data:    err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status, result, err := userHandler.UserService.LoginService(loginRequest, ctx)

	if err != nil {
		return c.Status(status).JSON(response.Output{
			Message: err.Error(),
			Time:    time.Now(),
			Data:    nil,
		})
	}

	return c.JSON(result)
}

func (userHandler *UserHandlerImpl) Register(c *fiber.Ctx) error {
	var registerRequest request.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if err := userHandler.Validator.Struct(registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Validation failed",
			Time:    time.Now(),
			Data:    err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	status, result, err := userHandler.UserService.RegisterService(registerRequest, ctx)

	if err != nil {
		return c.Status(status).JSON(response.Output{
			Message: err.Error(),
			Time:    time.Now(),
			Data:    nil,
		})
	}
	return c.JSON(result)
}
