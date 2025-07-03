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

	status, token, result := handler.UserService.LoginService(loginRequest, ctx)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",                           // Cookie for all path
		HTTPOnly: true,                          // Cookie can't be accessed from JavaScript
		Secure:   false,                         // Set to true for production via HTTPS
		SameSite: "None",                        // Important for cross-site cookies
		Domain:   "localhost",                   // Set to your domain in production
		Expires:  time.Now().Add(time.Hour * 1), // Set expired in 1 hour
	})

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

	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	return c.Status(status).JSON(result)
}
