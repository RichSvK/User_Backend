package test

import (
	"fmt"
	"net/http"
	"os"
	"stock_backend/internal/helper"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	requestBody := request.RegisterRequest{
		Email:    "test@gmail.com",
		Password: password,
		Username: "test_username_2",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/register"
	result, statusCode, err := PerformRequest[*response.RegisterResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "Registration successful", result.Message)
}

func TestRegisterBadRequest(t *testing.T) {
	requestBody := map[string]any{
		"Email":    "",
		"Password": 12345678,
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/register"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, domainerr.ErrInvalidRequestBody.Error(), result.Message)
}

func TestRegisterValidationFailed(t *testing.T) {
	requestBody := request.RegisterRequest{
		Email:    "",
		Password: "",
		Username: "",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/register"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "Email is required", result.Message)
}

func TestRegisterDuplicate(t *testing.T) {
	requestBody := request.RegisterRequest{
		Email:    email,
		Password: password,
		Username: "test_username",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/register"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, domainerr.ErrEmailExists.Error(), result.Message)
}

func TestLogin(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    email,
		Password: password,
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Login successful", result.Message)
	assert.NotEmpty(t, result.Token)
}

func TestLoginWrongPassword(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    email,
		Password: "12345678",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrWrongPassword.Error(), result.Message)
}

func TestLoginBadRequest(t *testing.T) {
	requestBody := map[string]any{
		"Email":    123456,
		"Password": "8765",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, domainerr.ErrInvalidRequestBody.Error(), result.Message)
}

func TestLoginMinPassword(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "test@gmail.com",
		Password: "8765",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "Password must be greater than or equal to 6", result.Message)
}

func TestLoginBadRequiredField(t *testing.T) {
	requestBody := request.LoginRequest{
		Password: "8765",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "Email is required", result.Message)
}

func TestLoginBadEmailField(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "test123",
		Password: "8765",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "Email must be a valid email address", result.Message)
}

func TestLoginNotFound(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "test_3@gmail.com",
		Password: password,
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, domainerr.ErrUserNotFound.Error(), result.Message)
}

func TestLogoutSuccess(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    email,
		Password: password,
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Login successful", result.Message)
	assert.NotEmpty(t, result.Token)

	delete(httpHeader, "Content-Type")
	httpHeader["Authorization"] = fmt.Sprintf("Bearer %s", result.Token)

	url = "/api/v1/auth/users/logout"
	logoutResult, statusCode, err := PerformRequest[*response.LogoutResponse](nil, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Logout successful", logoutResult.Message)
}

func TestLogoutFailed(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := "/api/v1/auth/users/logout"
	result, statusCode, err := PerformRequest[*response.LogoutResponse](nil, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestGetProfile(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    email,
		Password: password,
	}

	httpHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Login successful", result.Message)
	assert.NotEmpty(t, result.Token)

	delete(httpHeader, "Content-Type")
	httpHeader["Authorization"] = fmt.Sprintf("Bearer %s", result.Token)

	url = "/api/v1/auth/users/profile"
	profileResult, statusCode, err := PerformRequest[*response.UserProfileResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, email, profileResult.Email)
	assert.Equal(t, "test_username", profileResult.Username)
}

func TestGetProfileUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := "/api/v1/auth/users/profile"
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestVerifyUserTokenEmpty(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := "/api/v1/users/verify?token="
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, domainerr.ErrEmptyToken.Error(), result.Message)
}

func TestVerifyUserInvalidToken(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := "/api/v1/users/verify?token=123456789"
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrInvalidToken.Error(), result.Message)
}

func TestVerifyUser(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	var userId string
	var roleId int
	role := "user"
	err := db.QueryRow("SELECT id, roleId FROM users WHERE email = $1", email).
		Scan(&userId, &roleId)
	assert.Nil(t, err)

	if roleId == 2 {
		role = "admin"
	}

	verifyToken, err := helper.GenerateJWT(userId, email, role, os.Getenv("EMAIL_SECRET_KEY"))
	assert.Nil(t, err)

	url := "/api/v1/users/verify?token=" + verifyToken
	result, statusCode, err := PerformRequest[*response.VerifyResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "User verified successfully", result.Message)
}

func TestVerifyUserVerified(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	var userId string
	var roleId int
	role := "user"
	err := db.QueryRow("SELECT id, roleId FROM users WHERE email = $1", email).
		Scan(&userId, &roleId)
	assert.Nil(t, err)

	if roleId == 2 {
		role = "admin"
	}

	verifyToken, err := helper.GenerateJWT(userId, email, role, os.Getenv("EMAIL_SECRET_KEY"))
	assert.Nil(t, err)

	url := "/api/v1/users/verify?token=" + verifyToken
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, domainerr.ErrVerified.Error(), result.Message)
}

func TestDeleteUser(t *testing.T) {
	deletedEmail := "test_3@gmail.com"
	err := CreateTestUser(deletedEmail, password)
	assert.Nil(t, err)

	var userId string
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", deletedEmail).
		Scan(&userId)
	assert.Nil(t, err)

	requestBody := request.DeleteUserRequest{
		UserId: userId,
	}

	httpheader := map[string]string{
		"Authorization": "Bearer " + adminToken,
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}
	url := "/api/v1/auth/users"
	result, statusCode, err := PerformRequest[*response.DeleteUserResponse](requestBody, url, http.MethodDelete, httpheader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "User deleted successfully", result.Message)

	failedRes, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodDelete, httpheader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, domainerr.ErrUserNotFound.Error(), failedRes.Message)
}
