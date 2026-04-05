package test

import (
	"fmt"
	"net/http"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	watchlistPath    = "/api/v1/watchlists"
	addWatchlistPath = "/api/v1/watchlists/stocks"
)

func TestAddWatchlist(t *testing.T) {
	requestBody := request.AddWatchlistRequest{
		Stock: "NOBU",
	}

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	url := addWatchlistPath
	result, statusCode, err := PerformRequest[*response.AddWatchlistResponse](requestBody, url, http.MethodPost, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "Successfully added NOBU to watchlist", result.Message)
}

func TestAddWatchListUnauthorized(t *testing.T) {
	requestBody := request.AddWatchlistRequest{
		Stock: "NOBU",
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := addWatchlistPath
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestAddWatchlistDuplicate(t *testing.T) {
	requestBody := request.AddWatchlistRequest{
		Stock: "NOBU",
	}

	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	url := addWatchlistPath
	result, statusCode, err := PerformRequest[*response.FailedResponse](requestBody, url, http.MethodPost, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, domainerr.ErrWatchlistDuplicate.Error(), result.Message)
}

func TestAddWatchlistBadRequest(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := addWatchlistPath
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodPost, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, domainerr.ErrInvalidRequestBody.Error(), result.Message)
}

func TestGetWatchlist(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := watchlistPath
	result, statusCode, err := PerformRequest[*response.GetWatchlistResponse](nil, url, http.MethodGet, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Watchlist retrieved successfully", result.Message)
	assert.Equal(t, "NOBU", result.Stocks[0])
}

func TestGetWatchlistUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	url := watchlistPath
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestRemoveFromWatchlist(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	stock := "NOBU"
	url := fmt.Sprintf("%s/stocks/%s", watchlistPath, stock)
	result, statusCode, err := PerformRequest[*response.RemoveWatchlistResponse](nil, url, http.MethodDelete, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Successfully removed NOBU from watchlist", result.Message)
}

func TestRemoveFromWatchlistUnauthorized(t *testing.T) {
	httpHeader := map[string]string{
		"Accept": "application/json",
	}

	stock := "NOBU"
	url := fmt.Sprintf("%s/stocks/%s", watchlistPath, stock)
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodDelete, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, domainerr.ErrAuthorizationHeaderRequired.Error(), result.Message)
}

func TestRemoveFromWatchlistNotFound(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	stock := "NOBU"
	url := fmt.Sprintf("%s/stocks/%s", watchlistPath, stock)
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodDelete, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, domainerr.ErrWatchlistNotFound.Error(), result.Message)
}

func TestGetWatchlistNotFound(t *testing.T) {
	httpHeader := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
	}

	url := watchlistPath
	result, statusCode, err := PerformRequest[*response.FailedResponse](nil, url, http.MethodGet, httpHeader)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, domainerr.ErrWatchlistNotFound.Error(), result.Message)
}
