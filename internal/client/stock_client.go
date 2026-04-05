package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/response"
	"strings"
	"time"

	"github.com/sony/gobreaker"
)

type StockClient interface {
	GetStock(ctx context.Context, ticker string) error
}

type stockClient struct {
	httpClient *http.Client
	breaker    *gobreaker.CircuitBreaker
	baseUrl    string
}

func NewStockClient(url string, breaker *gobreaker.CircuitBreaker) StockClient {
	return &stockClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		breaker:    breaker,
		baseUrl:    url,
	}
}

func (c *stockClient) GetStock(ctx context.Context, stock string) error {
	_, err := c.breaker.Execute(func() (interface{}, error) {
		url := fmt.Sprintf("%s/api/v1/stocks?code=%s", c.baseUrl, stock)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		// Call the request
		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		// If the response is not ok
		if res.StatusCode != http.StatusOK {
			var webResponse response.WebResponse
			if err := json.NewDecoder(res.Body).Decode(&webResponse); err != nil {
				return nil, err
			}

			err := domainerr.NewServiceError(res.StatusCode, webResponse.Message)
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return domainerr.ErrStockServiceUnavailable
		}

		return err
	}

	return nil
}

func mapClientError(err error) error {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return domainerr.ErrStockServiceUnavailable
	}

	if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
		return domainerr.ErrServiceTimeout
	}
	return domainerr.ErrStockServiceUnavailable
}
