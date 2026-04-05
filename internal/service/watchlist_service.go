package service

import (
	"context"
	"fmt"
	"stock_backend/internal/client"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/response"
	"stock_backend/internal/repository"
)

type WatchlistService interface {
	AddToWatchlist(ctx context.Context, userId string, stock string) (*response.AddWatchlistResponse, error)
	RemoveFromWatchlist(ctx context.Context, userId string, stock string) (*response.RemoveWatchlistResponse, error)
	GetWatchlist(ctx context.Context, userId string) (*response.GetWatchlistResponse, error)
}

type WatchlistServiceImpl struct {
	Repository repository.WatchlistRepository
	stockClient client.StockClient
}

func NewWatchlistService(repository repository.WatchlistRepository, stockClient client.StockClient) WatchlistService {
	return &WatchlistServiceImpl{
		Repository: repository,
		stockClient: stockClient,
	}
}

func (service *WatchlistServiceImpl) AddToWatchlist(ctx context.Context, userId string, stock string) (*response.AddWatchlistResponse, error) {
	if err := service.stockClient.GetStock(ctx, stock); err != nil {
		return nil, err
	}

	err := service.Repository.AddWatchlist(ctx, userId, stock)
	if err != nil {
		return nil, err
	}

	response := &response.AddWatchlistResponse{
		Message: fmt.Sprintf("Successfully added %s to watchlist", stock),
	}

	return response, nil
}

func (service *WatchlistServiceImpl) RemoveFromWatchlist(ctx context.Context, userId string, stock string) (*response.RemoveWatchlistResponse, error) {
	if err := service.Repository.RemoveWatchlist(ctx, userId, stock); err != nil {
		return nil, err
	}

	response := &response.RemoveWatchlistResponse{
		Message: fmt.Sprintf("Successfully removed %s from watchlist", stock),
	}
	return response, nil
}

func (service *WatchlistServiceImpl) GetWatchlist(ctx context.Context, userId string) (*response.GetWatchlistResponse, error) {
	watchlist, err := service.Repository.GetWatchlistByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(watchlist) == 0 {
		return nil, domainerr.ErrWatchlistNotFound
	}

	response := &response.GetWatchlistResponse{
		Message: "Watchlist retrieved successfully",
		Stocks:  watchlist,
	}
	return response, nil
}
