package handler

import (
	"context"
	"fmt"
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
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	if err := handler.Validator.Struct(loginRequest); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
			firstErr := errs[0]
			field := firstErr.Field()
			tag := firstErr.Tag()

			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("%s is required", field)
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", field, firstErr.Param())
			case "email":
				msg = fmt.Sprintf("%s must be a valid email address", field)
			default:
				msg = fmt.Sprintf("%s is invalid", field)
			}

			return c.Status(fiber.StatusBadRequest).JSON(response.Output{
				Message: msg,
				Time:    time.Now(),
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "Invalid request failed",
			Time:    time.Now(),
			Data:    nil,
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

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "User ID is required",
			Time:    time.Now(),
			Data:    nil,
		})
	}
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

func (handler *UserHandlerImpl) GetUserInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Output{
			Message: "User ID is required",
			Time:    time.Now(),
			Data:    nil,
		})
	}

	status, result := handler.UserService.GetUserProfile(userId, ctx)

	return c.Status(status).JSON(result)
}
