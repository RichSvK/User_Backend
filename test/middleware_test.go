package test

import (
	"net/http"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRole(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "CC",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "Add CC to favorite success", result.Message)
}

func TestUserRoleUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + adminToken,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "CC",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusForbidden, statusCode)
	assert.Equal(t, "Unauthorized access", result.Message)
}

func TestAdminUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := "/api/v1/auth/users"
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusForbidden, statusCode)
	assert.Equal(t, "Unauthorized access", result.Message)
}

func TestLoggedOutFailed(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := "/api/v1/users/login"
	result, statusCode, err := PerformRequest[*response.LoginResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "You are already logged in", result.Message)
}
