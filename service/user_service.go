package service

import (
	"context"
	"stock_backend/helper"
	"stock_backend/model/entity"
	"stock_backend/model/request"
	"stock_backend/model/response"
	"stock_backend/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	LoginService(request request.LoginRequest, ctx context.Context) (int, any)
	RegisterService(request request.RegisterRequest, ctx context.Context) (int, any)
	LogOutService(userId string, ctx context.Context) (int, any)
	DeleteUserService(userId string, ctx context.Context) (int, any)
	GetUserProfile(userId string, ctx context.Context) (int, any)
}

type UserServiceImpl struct {
	Repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceImpl{
		Repository: repository,
	}
}

func (service *UserServiceImpl) LoginService(request request.LoginRequest, ctx context.Context) (int, any) {
	user, err := service.Repository.GetUser(request.Email, ctx)
	if err != nil {
		return fiber.StatusNotFound,
			response.Output{
				Message: "User not found",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return fiber.StatusUnauthorized,
			response.Output{
				Message: "Wrong password",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	userId := user.ID.String()
	token, err := helper.GenerateJWT(userId, user.Email, user.Role)
	if err != nil {
		return fiber.StatusInternalServerError,
			response.Output{
				Message: "Internal server error",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	return fiber.StatusOK,
		response.Output{
			Message: "Login Success",
			Time:    time.Now(),
			Data: map[string]string{
				"token": token,
			},
		}
}

func (service *UserServiceImpl) RegisterService(request request.RegisterRequest, ctx context.Context) (int, any) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.StatusInternalServerError,
			response.Output{
				Message: "Internal server error",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	user := entity.User{
		ID:       uuid.New(),
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	if err := service.Repository.Create(user, ctx); err != nil {
		return fiber.StatusInternalServerError,
			response.Output{
				Message: "Internal server error",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	return fiber.StatusOK,
		response.Output{
			Message: "Register Success",
			Time:    time.Now(),
			Data:    nil,
		}
}

func (service *UserServiceImpl) LogOutService(userId string, ctx context.Context) (int, any) {
	if err := service.Repository.Logout(userId, ctx); err != nil {
		return fiber.StatusInternalServerError, response.Output{
			Message: "Internal server error",
			Time:    time.Now(),
			Data:    nil,
		}
	}

	return fiber.StatusOK, response.Output{
		Message: "Logout Success",
		Time:    time.Now(),
		Data:    nil,
	}
}

func (service *UserServiceImpl) DeleteUserService(userId string, ctx context.Context) (int, any) {
	if err := service.Repository.DeleteUser(userId, ctx); err != nil {
		return fiber.StatusInternalServerError, response.Output{
			Message: "Internal server error",
			Time:    time.Now(),
			Data:    nil,
		}
	}

	return fiber.StatusOK, response.Output{
		Message: "User deleted successfully",
		Time:    time.Now(),
		Data:    nil,
	}
}

func (service *UserServiceImpl) GetUserProfile(userId string, ctx context.Context) (int, any) {
	user, err := service.Repository.GetUserByID(userId, ctx)
	if err != nil {
		return fiber.StatusNotFound,
			response.Output{
				Message: "User not found",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	return fiber.StatusOK,
		response.Output{
			Message: "User profile retrieved successfully",
			Time:    time.Now(),
			Data: map[string]string{
				"username": user.Username,
				"email":    user.Email,
			},
		}
}
