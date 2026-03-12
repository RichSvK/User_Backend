package service

import (
	"context"
	"fmt"
	"log"
	"stock_backend/internal/helper"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/entity"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"stock_backend/internal/repository"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(request request.LoginRequest, ctx context.Context) (*response.LoginResponse, error)
	Register(request request.RegisterRequest, ctx context.Context) (*response.RegisterResponse, error)
	VerifyUser(tokenString string, ctx context.Context) (*response.VerifyResponse, error)
	Logout(userId string, ctx context.Context) (*response.LogoutResponse, error)
	DeleteUser(userId string, ctx context.Context) (*response.DeleteUserResponse, error)
	GetProfile(userId string, ctx context.Context) (*response.UserProfileResponse, error)
}

type UserServiceImpl struct {
	Repository repository.UserRepository
	JwtSecret  string
	Smtp       smtpConfig
}

func NewUserService(repository repository.UserRepository, jwtSecret string, smtp smtpConfig) UserService {
	return &UserServiceImpl{
		Repository: repository,
		JwtSecret:  jwtSecret,
		Smtp:       smtp,
	}
}

func (service *UserServiceImpl) Login(request request.LoginRequest, ctx context.Context) (*response.LoginResponse, error) {
	user, err := service.Repository.GetUser(request.Email, ctx)
	if err != nil {
		return nil, err
	}

	// Run this code if email verification is required
	// if !user.Verified {
	// 	if err := SendVerificationEmail(ctx, *user); err != nil {
	// 		log.Println("email failed:", err)
	// 	}

	// 	return nil, domainerr.ErrNotVerified
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, domainerr.ErrWrongPassword
	}

	userId := user.ID.String()

	token, err := helper.GenerateJWT(userId, user.Email, user.Role, service.JwtSecret)
	if err != nil {
		return nil, domainerr.ErrInternal
	}

	response := &response.LoginResponse{
		Message: "Login successful",
		Token:   token,
	}

	return response, nil
}

func (service *UserServiceImpl) Register(request request.RegisterRequest, ctx context.Context) (*response.RegisterResponse, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, domainerr.ErrInternal
	}

	user := entity.User{
		ID:       uuid.New(),
		Username: request.Username,
		Email:    request.Email,
		Password: string(hash),
	}

	// var createdUser *entity.User
	if _, err := service.Repository.Create(user, ctx); err != nil {
		return nil, err
	}

	token, err := helper.GenerateJWT(user.ID.String(), user.Email, user.Role, service.Smtp.Secret)
	if err != nil {
		return nil, err
	}

	verifyURL := fmt.Sprintf("http://%s:%s/api/v1/users/verify?token=%s", service.Smtp.AppHost, service.Smtp.AppPort, token)
	htmlBody, err := renderVerificationEmail(verifyURL)
	if err != nil {
		return nil, err
	}

	// Run this code if email verification is required and SMTP server is configured
	if err := service.Smtp.sendHTML(ctx, user.Email, "Verify your Stock App account", htmlBody); err != nil {
		log.Println("email failed:", err)
	}

	response := &response.RegisterResponse{
		Message: "Registration successful",
	}

	return response, nil
}

func (service *UserServiceImpl) VerifyUser(tokenString string, ctx context.Context) (*response.VerifyResponse, error) {
	token, err := helper.ValidateJWT(tokenString, service.Smtp.Secret)
	if err != nil {
		return nil, domainerr.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, domainerr.ErrInvalidTokenClaims
	}

	userId, ok := claims["sub"].(string)
	if !ok || userId == "" {
		return nil, domainerr.ErrInvalidTokenClaims
	}

	if err := service.Repository.VerifyUser(userId, ctx); err != nil {
		return nil, err
	}

	response := &response.VerifyResponse{
		Message: "User verified successfully",
	}

	return response, nil
}

func (service *UserServiceImpl) Logout(userId string, ctx context.Context) (*response.LogoutResponse, error) {
	if err := service.Repository.Logout(userId, ctx); err != nil {
		return nil, err
	}

	response := &response.LogoutResponse{
		Message: "Logout successful",
	}

	return response, nil
}

func (service *UserServiceImpl) DeleteUser(userId string, ctx context.Context) (*response.DeleteUserResponse, error) {
	if err := service.Repository.DeleteUser(userId, ctx); err != nil {
		return nil, err
	}

	response := &response.DeleteUserResponse{
		Message: "User deleted successfully",
	}

	return response, nil
}

func (service *UserServiceImpl) GetProfile(userId string, ctx context.Context) (*response.UserProfileResponse, error) {
	user, err := service.Repository.GetUserByID(userId, ctx)

	if err != nil {
		return nil, err
	}

	response := &response.UserProfileResponse{
		Username: user.Username,
		Email:    user.Email,
	}
	return response, nil
}
