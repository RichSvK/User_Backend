package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"stock_backend/helper"
	"stock_backend/model/entity"
	domain_error "stock_backend/model/error"
	"stock_backend/model/request"
	"stock_backend/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(request request.LoginRequest, ctx context.Context) (string, error)
	Register(request request.RegisterRequest, ctx context.Context) error
	VerifyUser(tokenString string, ctx context.Context) error
	Logout(userId string, ctx context.Context) error
	DeleteUser(userId string, ctx context.Context) error
	GetProfile(userId string, ctx context.Context) (map[string]string, error)
}

type UserServiceImpl struct {
	Repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceImpl{
		Repository: repository,
	}
}

func (service *UserServiceImpl) Login(request request.LoginRequest, ctx context.Context) (string, error) {
	user, err := service.Repository.GetUser(request.Email, ctx)
	if err != nil {
		return "", err
	}

	if !user.Verified {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := SendVerificationEmail(ctx, *user); err != nil {
				log.Println("email failed:", err)
			}
		}()

		return "", domain_error.ErrNotVerified
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return "", domain_error.ErrWrongPassword
	}

	userId := user.ID.String()

	return helper.GenerateJWT(userId, user.Email, user.Role)
}

func (service *UserServiceImpl) Register(request request.RegisterRequest, ctx context.Context) error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return domain_error.ErrInternal
	}

	user := entity.User{
		ID:       uuid.New(),
		Username: request.Username,
		Email:    request.Email,
		Password: string(hash),
	}

	var createdUser *entity.User
	if createdUser, err = service.Repository.Create(user, ctx); err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := SendVerificationEmail(ctx, *createdUser); err != nil {
			log.Println("email failed:", err)
		}
	}()

	return nil
}

func (service *UserServiceImpl) VerifyUser(tokenString string, ctx context.Context) error {
	token, err := helper.ValidateJWT(tokenString)
	if err != nil {
		return domain_error.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain_error.ErrInvalidTokenClaims
	}

	userId := claims["sub"].(string)

	if err := service.Repository.VerifyUser(userId, ctx); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) Logout(userId string, ctx context.Context) error {
	if err := service.Repository.Logout(userId, ctx); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) DeleteUser(userId string, ctx context.Context) error {
	if err := service.Repository.DeleteUser(userId, ctx); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) GetProfile(userId string, ctx context.Context) (map[string]string, error) {
	user, err := service.Repository.GetUserByID(userId, ctx)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}, nil
}

func SendVerificationEmail(ctx context.Context, user entity.User) error {
	smtpUser := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	appHost := os.Getenv("APP_HOST")

	if smtpUser == "" || smtpPass == "" || smtpHost == "" || smtpPort == "" || appHost == "" {
		return errors.New("missing smtp configuration")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	token, err := helper.GenerateJWT(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return err
	}

	verifyURL := fmt.Sprintf(
		"%s/api/user/verify?token=%s",
		appHost,
		token,
	)

	subject := "Verify your Stock App account"

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

	msg := []byte(
		"From: " + smtpUser + "\r\n" +
			"To: " + user.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
			htmlBody,
	)

	// respect context cancellation
	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(
			smtpHost+":"+smtpPort,
			auth,
			smtpUser,
			[]string{user.Email},
			msg,
		)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
