package service

import (
	"context"
	"fmt"
	"stock_backend/model/entity"
	"stock_backend/model/response"
	"stock_backend/repository"
)

type FavoriteService interface {
	CreateFavorite(userId string, underwriterId string, ctx context.Context) (*response.AddFavorite, error)
	GetFavorites(userId string, ctx context.Context) (*response.GetFavorites, error)
	RemoveFavorite(userId string, underwriterCode string, ctx context.Context) (*response.RemoveFavorite, error)
}

type FavoriteServiceImpl struct {
	Repository repository.FavoriteRepository
}

func NewFavoriteService(repository repository.FavoriteRepository) FavoriteService {
	return &FavoriteServiceImpl{
		Repository: repository,
	}
}

func (service *FavoriteServiceImpl) CreateFavorite(userId string, underwriterId string, ctx context.Context) (*response.AddFavorite, error) {
	favorite := &entity.Favorite{
		UserID:        userId,
		UnderwriterID: underwriterId,
	}

	if err := service.Repository.Create(favorite, ctx); err != nil {
		return nil, err
	}

	response := &response.AddFavorite{
		Message: fmt.Sprintf("Add %s to favorite success", underwriterId),
	}

	return response, nil
}

func (service *FavoriteServiceImpl) GetFavorites(userId string, ctx context.Context) (*response.GetFavorites, error) {
	favoriteData, err := service.Repository.GetFavorites(userId, ctx)
	if err != nil {
		return nil, err
	}

	response := &response.GetFavorites{
		Message: "Favorite Found",
		Data:    favoriteData,
	}
	return response, nil
}

func (service *FavoriteServiceImpl) RemoveFavorite(userId string, underwriterCode string, ctx context.Context) (*response.RemoveFavorite, error) {
	if err := service.Repository.RemoveFavorite(userId, underwriterCode, ctx); err != nil {
		return nil, err
	}

	response := &response.RemoveFavorite{
		Message: "Remove favorite success",
	}

	return response, nil
}
