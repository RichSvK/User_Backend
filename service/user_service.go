package service

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"stock_backend/helper"
	"stock_backend/model/entity"
	"stock_backend/model/request"
	"stock_backend/model/response"
	"stock_backend/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	LoginService(request request.LoginRequest, ctx context.Context) (int, any)
	RegisterService(request request.RegisterRequest, ctx context.Context) (int, any)
	VerifyUserService(tokenString string, ctx context.Context) (int, any)
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

	var createdUser *entity.User
	if createdUser, err = service.Repository.Create(user, ctx); err != nil {
		return fiber.StatusInternalServerError,
			response.Output{
				Message: "Error registering user",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	go func() {
		SendVerificationEmail(*createdUser)
	}()

	return fiber.StatusOK,
		response.Output{
			Message: "Register Success",
			Time:    time.Now(),
			Data:    nil,
		}
}

func (service *UserServiceImpl) VerifyUserService(tokenString string, ctx context.Context) (int, any) {
	token, err := helper.ValidateJWT(tokenString)
	if err != nil {
		return fiber.StatusBadRequest,
			response.Output{
				Message: "Invalid or expired token",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fiber.StatusBadRequest,
			response.Output{
				Message: "Invalid token claims",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	userId := claims["sub"].(string)

	if err := service.Repository.VerifyUser(userId, ctx); err != nil {
		return fiber.StatusInternalServerError,
			response.Output{
				Message: "Internal server error",
				Time:    time.Now(),
				Data:    nil,
			}
	}

	return fiber.StatusOK,
		response.Output{
			Message: "Account verified successfully",
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

func SendVerificationEmail(user entity.User) {
	smtpUser := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPass,
		smtpHost,
	)

	token, err := helper.GenerateJWT(user.ID.String(), user.Email, user.Role)
	if err != nil {
		log.Println("Failed to generate verification token:", err)
		return
	}

	appHost := os.Getenv("APP_HOST")
	verifyURL := fmt.Sprintf(
		"%s/api/user/verify?token=%s",
		appHost,
		token,
	)

	subject := "Verify your Stock App account"

	// HTML email body with button
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
		<table width="100%%" cellpadding="0" cellspacing="0">
			<tr>
				<td align="center">
					<table width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; padding: 30px; border-radius: 8px;">
						<tr>
							<td align="center">
								<h2>Verify your email address</h2>
								<p>Thanks for signing up! Please confirm your email address by clicking the button below.</p>
								<a href="%s"
								style="
									display: inline-block;
									padding: 14px 24px;
									margin-top: 20px;
									background-color: #007bff;
									color: #ffffff;
									text-decoration: none;
									border-radius: 6px;
									font-weight: bold;
								">
									Verify Account
								</a>
								<p style="margin-top: 30px; font-size: 12px; color: #777;">
									If you didn’t create an account, you can safely ignore this email.
								</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, verifyURL)

	// Email headers (VERY IMPORTANT)
	msg := []byte(
		"From: " + smtpUser + "\r\n" +
			"To: " + user.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			htmlBody,
	)

	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpUser,
		[]string{user.Email},
		msg,
	)

	if err != nil {
		log.Println("Failed to send verification email:", err)
	}
}
