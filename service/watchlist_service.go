package service

import (
	"context"
	"fmt"
	domain_error "stock_backend/model/error"
	"stock_backend/repository"
)

type WatchlistService interface {
	AddToWatchlist(ctx context.Context, userId string, stock string) error
	RemoveFromWatchlist(ctx context.Context, userId string, stock string) error
	GetWatchlist(ctx context.Context, userId string) ([]string, error)
}

type WatchlistServiceImpl struct {
	Repository repository.WatchlistRepository
}

func NewWatchlistService(repository repository.WatchlistRepository) WatchlistService {
	return &WatchlistServiceImpl{
		Repository: repository,
	}
}

func (service *WatchlistServiceImpl) AddToWatchlist(ctx context.Context, userId string, stock string) error {
	fmt.Println("Adding to watchlist:", userId, stock)
	err := service.Repository.AddWatchlist(ctx, userId, stock)
	return err
}

func (service *WatchlistServiceImpl) RemoveFromWatchlist(ctx context.Context, userId string, stock string) error {
	err := service.Repository.RemoveWatchlist(ctx, userId, stock)
	return err
}

func (service *WatchlistServiceImpl) GetWatchlist(ctx context.Context, userId string) ([]string, error) {
	watchlist, err := service.Repository.GetWatchlistByUserID(ctx, userId)
	fmt.Println("Adding to watchlist:", userId)

	if err != nil {
		return nil, err
	}

	if len(watchlist) == 0 {
		return nil, domain_error.ErrWatchlistNotFound
	}

	return watchlist, nil
}
