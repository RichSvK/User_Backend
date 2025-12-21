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
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Invalid request",
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

			return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
				Message: msg,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Invalid request",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := handler.UserService.Login(loginRequest, ctx)
	if err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"message": "Login successful",
		"token":   result,
	})
}

func (handler *UserHandlerImpl) Register(c *fiber.Ctx) error {
	var registerRequest request.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Invalid request",
		})
	}

	if err := handler.Validator.Struct(registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Validation failed",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := handler.UserService.Register(registerRequest, ctx); err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.Output{
		Message: "User registered successfully",
	})
}

func (handler *UserHandlerImpl) VerifyUser(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Token is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := handler.UserService.VerifyUser(token, ctx); err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Output{
		Message: "User verified successfully",
	})
}

func (handler *UserHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "User ID is required",
		})
	}

	if err := handler.UserService.Logout(userId, ctx); err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Output{
		Message: "Logout successful",
	})
}

func (handler *UserHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var deleteRequest request.DeleteUserRequest
	if err := c.BodyParser(&deleteRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Invalid request",
		})
	}

	if err := handler.Validator.Struct(deleteRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "Validation failed",
		})
	}

	if err := handler.UserService.DeleteUser(deleteRequest.UserId, ctx); err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Output{
		Message: "User deleted successfully",
	})
}

func (handler *UserHandlerImpl) GetUserInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userId := c.Get("X-User-ID")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.FailedResponse{
			Message: "User ID is required",
		})
	}

	result, err := handler.UserService.GetProfile(userId, ctx)
	if err != nil {
		return c.Status(MapErrorToHTTPStatus(err)).JSON(response.FailedResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
