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
	Logout(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
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
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if err := handler.Validator.Struct(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Validation failed",
			Time:    time.Now(),
			Data:    err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status, result := handler.UserService.LoginService(loginRequest, ctx)

	return c.Status(status).JSON(result)
}

func (handler *UserHandlerImpl) Register(c *fiber.Ctx) error {
	var registerRequest request.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if err := handler.Validator.Struct(registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Validation failed",
			Time:    time.Now(),
			Data:    err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	status, result := handler.UserService.RegisterService(registerRequest, ctx)

	return c.Status(status).JSON(result)
}

func (handler *UserHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId := c.Locals("userId").(string)
	status, result := handler.UserService.LogOutService(userId, ctx)

	return c.Status(status).JSON(result)
}

func (handler *UserHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var deleteRequest request.DeleteUserRequest
	if err := c.BodyParser(&deleteRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if err := handler.Validator.Struct(deleteRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Validation failed",
			Time:    time.Now(),
			Data:    err.Error(),
		})
	}

	status, result := handler.UserService.DeleteUserService(deleteRequest.UserId, ctx)

	return c.Status(status).JSON(result)
}
