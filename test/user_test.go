package test

import (
	"fmt"
	"net/http"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ClearUser()
	requestBody := request.RegisterRequest{
		Email:    "richardsugiharto0@gmail.com",
		Password: "87654321",
		Username: "rich_svk",
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
	result, statusCode, err := PerformRequest[*response.RegisterResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "Validation failed", result.Message)
}

func TestRegisterDuplicate(t *testing.T) {
	requestBody := request.RegisterRequest{
		Email:    "richardsugiharto0@gmail.com",
		Password: "87654321",
		Username: "rich_svk",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/register"
	result, statusCode, err := PerformRequest[*response.RegisterResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "internal server error", result.Message)
}

func TestLogin(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "richardsugiharto0@gmail.com",
		Password: "87654321",
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
		Email:    "richardsugiharto0@gmail.com",
		Password: "12345678",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "wrong password", result.Message)
}

func TestLoginBadRequest(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "richardsugiharto0@gmail.com",
		Password: "8765",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "Password must be at least 6 characters", result.Message)
}

func TestLogoutSuccess(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "richardsugiharto0@gmail.com",
		Password: "87654321",
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
	assert.Equal(t, "Authorization header is required", result.Message)
}

func TestGetProfile(t *testing.T) {
	requestBody := request.LoginRequest{
		Email:    "richardsugiharto0@gmail.com",
		Password: "87654321",
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
	assert.Equal(t, "richardsugiharto0@gmail.com", profileResult.Email)
	assert.Equal(t, "rich_svk", profileResult.Username)
}

func TestGetProfileFailed(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := "/api/v1/auth/users/profile"
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "Authorization header is required", result.Message)
}
