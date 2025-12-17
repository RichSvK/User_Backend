package service

import (
	"context"
	"fmt"
	"stock_backend/model/response"
	"stock_backend/repository"
)

type WatchlistService interface {
	AddToWatchlist(ctx context.Context, userId string, stock string) (int, any)
	RemoveFromWatchlist(ctx context.Context, userId string, stock string) (int, any)
	GetWatchlist(ctx context.Context, userId string) (int, any)
}

type WatchlistServiceImpl struct {
	Repository repository.WatchlistRepository
}

func NewWatchlistService(repository repository.WatchlistRepository) WatchlistService {
	return &WatchlistServiceImpl{
		Repository: repository,
	}
}
func (service *WatchlistServiceImpl) AddToWatchlist(ctx context.Context, userId string, stock string) (int, any) {
	err := service.Repository.AddWatchlist(ctx, userId, stock)

	if err != nil {
		return 500, response.WatchlistResponse{
			Message: fmt.Sprintf("Failed to add %s to watchlist", stock),
		}
	}

	return 201, response.WatchlistResponse{
		Message: fmt.Sprintf("Successfully added %s to watchlist", stock),
	}
}

func (service *WatchlistServiceImpl) RemoveFromWatchlist(ctx context.Context, userId string, stock string) (int, any) {
	err := service.Repository.RemoveWatchlist(ctx, userId, stock)

	if err != nil {
		return 500, response.WatchlistResponse{
			Message: fmt.Sprintf("Failed to remove %s from watchlist", stock),
		}
	}

	return 200, response.WatchlistResponse{
		Message: fmt.Sprintf("Successfully removed %s from watchlist", stock),
	}
}

func (service *WatchlistServiceImpl) GetWatchlist(ctx context.Context, userId string) (int, any) {
	watchlist, err := service.Repository.GetWatchlistByUserID(ctx, userId)
	if err != nil {
		return 500, response.GetWatchlistResponse{
			Message: "Failed to retrieve watchlist",
			Stocks:  nil,
		}
	}

	if len(watchlist) == 0 {
		return 200, response.GetWatchlistResponse{
			Message: "User doesn't have any watchlist",
			Stocks:  nil,
		}
	}

	return 200, response.GetWatchlistResponse{
		Message: "Watchlist retrieved successfully",
		Stocks:  watchlist,
	}
}
