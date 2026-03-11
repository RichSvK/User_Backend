package test

import (
	"net/http"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

const userEmail = "richardsugiharto0@gmail.com"
const password = "87654321"

func TestAddFavorites(t *testing.T) {
	ClearTable("favorites")
	token, err := GetUserToken(userEmail, password)
	assert.Nil(t, err)

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "KI",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "Add KI to favorite success", result.Message)
}

func TestAddFavoritesDuplicates(t *testing.T) {
	token, err := GetUserToken(userEmail, password)
	assert.Nil(t, err)

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "KI",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, "Data already exists", result.Message)
}

func TestAddFavoritesUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "KI",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "Authorization header is required", result.Message)
}

func TestGetFavorites(t *testing.T) {
	token, err := GetUserToken(userEmail, password)
	assert.Nil(t, err)

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.GetFavoritesResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Favorite Found", result.Message)
	assert.Equal(t, "KI", result.Data[0])
}

func TestGetFavoritesUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := "/api/v1/auth/favorites"
	result, statusCode, err := PerformRequest[*response.GetFavoritesResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "Authorization header is required", result.Message)
}

func TestRemoveFavoriteResponses(t *testing.T) {
	token, err := GetUserToken(userEmail, password)
	assert.Nil(t, err)

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	underwriterId := "KI"
	url := "/api/v1/auth/favorites/" + underwriterId

	result, statusCode, err := PerformRequest[*response.RemoveFavoriteResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Remove favorite success", result.Message)
}

func TestRemoveFavoriteResponsesUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	underwriterId := "KI"
	url := "/api/v1/auth/favorites/" + underwriterId

	result, statusCode, err := PerformRequest[*response.RemoveFavoriteResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "Authorization header is required", result.Message)
}

func TestRemoveFavoriteResponsesNotFound(t *testing.T) {
	token, err := GetUserToken(userEmail, password)
	assert.Nil(t, err)

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept": "application/json",
	}

	underwriterId := "KI"
	url := "/api/v1/auth/favorites/" + underwriterId

	result, statusCode, err := PerformRequest[*response.RemoveFavoriteResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, "Favorites not found", result.Message)
}
