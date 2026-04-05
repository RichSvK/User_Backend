package service

import (
	"context"
	"fmt"
	"stock_backend/internal/entity"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/response"
	"stock_backend/internal/repository"
)

type FavoriteService interface {
	CreateFavorite(ctx context.Context, userId string, underwriterId string) (*response.AddFavoriteResponse, error)
	GetFavorites(ctx context.Context, userId string) (*response.GetFavoritesResponse, error)
	RemoveFavorite(ctx context.Context, userId string, underwriterCode string) (*response.RemoveFavoriteResponse, error)
}

type FavoriteServiceImpl struct {
	Repository repository.FavoriteRepository
}

func NewFavoriteService(repository repository.FavoriteRepository) FavoriteService {
	return &FavoriteServiceImpl{
		Repository: repository,
	}
}

func (service *FavoriteServiceImpl) CreateFavorite(ctx context.Context, userId string, underwriterId string) (*response.AddFavoriteResponse, error) {
	favorite := &entity.Favorite{
		UserID:        userId,
		UnderwriterID: underwriterId,
	}

	if err := service.Repository.Create(favorite, ctx); err != nil {
		return nil, err
	}

	response := &response.AddFavoriteResponse{
		Message: fmt.Sprintf("Add %s to favorite success", underwriterId),
	}

	return response, nil
}

func (service *FavoriteServiceImpl) GetFavorites(ctx context.Context, userId string) (*response.GetFavoritesResponse, error) {
	favoriteData, err := service.Repository.GetFavorites(userId, ctx)
	if err != nil {
		return nil, err
	}

	if len(favoriteData) == 0 {
		return nil, domainerr.ErrFavoritesNotFound
	}

	response := &response.GetFavoritesResponse{
		Message: "Favorite Found",
		Data:    favoriteData,
	}
	return response, nil
}

func (service *FavoriteServiceImpl) RemoveFavorite(ctx context.Context, userId string, underwriterCode string) (*response.RemoveFavoriteResponse, error) {
	if err := service.Repository.RemoveFavorite(userId, underwriterCode, ctx); err != nil {
		return nil, err
	}

	response := &response.RemoveFavoriteResponse{
		Message: "Remove favorite success",
	}

	return response, nil
}
