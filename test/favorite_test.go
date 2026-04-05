package test

import (
	"fmt"
	"net/http"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	favoritesPath = "/api/v1/favorites"
)

func TestAddFavorites(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "KI",
	}

	url := favoritesPath
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "Add KI to favorite success", result.Message)
}

func TestAddFavoritesDuplicates(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "KI",
	}

	url := favoritesPath
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, domainerr.ErrFavoritesDuplicate.Error(), result.Message)
}

func TestAddFavoritesUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := request.AddFavoriteUnderwriterRequest{
		UnderwriterId: "KI",
	}

	url := favoritesPath
	result, statusCode, err := PerformRequest[*response.AddFavoriteResponse](requestBody, url, http.MethodPost, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestGetFavorites(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := favoritesPath
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

	url := favoritesPath
	result, statusCode, err := PerformRequest[*response.GetFavoritesResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestRemoveFavorite(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	underwriterId := "KI"
	url := fmt.Sprintf("%s/%s", favoritesPath, underwriterId)
	result, statusCode, err := PerformRequest[*response.RemoveFavoriteResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Remove favorite success", result.Message)
}

func TestRemoveFavoriteUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	underwriterId := "KI"
	url := fmt.Sprintf("%s/%s", favoritesPath, underwriterId)
	result, statusCode, err := PerformRequest[*response.RemoveFavoriteResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestRemoveFavoriteBadRequest(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := fmt.Sprintf("%s/A321", favoritesPath)
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, domainerr.ErrFavoritesUnderwriterInvalid.Error(), result.Message)
}

func TestRemoveFavoriteNotFound(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	underwriterId := "KI"
	url := fmt.Sprintf("%s/%s", favoritesPath, underwriterId)
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodDelete, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, domainerr.ErrFavoritesNotFound.Error(), result.Message)
}

func TestGetFavoritesNotFound(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := favoritesPath
	result, statusCode, err := PerformRequest[*response.GetFavoritesResponse](nil, url, http.MethodGet, httpHeader)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, domainerr.ErrFavoritesNotFound.Error(), result.Message)
}
