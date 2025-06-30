package service

import (
	"context"
	"stock_backend/model/entity"
	"stock_backend/model/request"
	"stock_backend/model/response"
	"stock_backend/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteService interface {
	CreateFavorite(userId string, request request.FavoriteUnderwriterRequest, ctx context.Context) (int, any)
	GetFavoritesService(userId string, ctx context.Context) (int, any)
}

type FavoriteServiceImpl struct {
	Repository repository.FavoriteRepository
}

func NewFavoriteService(repository repository.FavoriteRepository) FavoriteService {
	return &FavoriteServiceImpl{
		Repository: repository,
	}
}

func (s *FavoriteServiceImpl) CreateFavorite(userId string, request request.FavoriteUnderwriterRequest, ctx context.Context) (int, any) {
	favorite := &entity.Favorite{
		ID:           primitive.NewObjectID(),
		UserID:       userId,
		Underwriters: request.Underwriter,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.Repository.Create(favorite, ctx)
	if err != nil {
		result := response.Output{
			Message: "Add favorite failed",
			Time:    time.Now(),
			Data:    nil,
		}
		return fiber.StatusBadRequest, result
	}

	result := response.Output{
		Message: "Add favorite success",
		Time:    time.Now(),
		Data:    nil,
	}
	return fiber.StatusOK, result
}

func (s *FavoriteServiceImpl) GetFavoritesService(userId string, ctx context.Context) (int, any) {
	favoriteData, err := s.Repository.GetFavorites(userId, ctx)
	if err != nil {
		result := response.Output{
			Message: err.Error(),
			Time:    time.Now(),
			Data:    nil,
		}
		return fiber.StatusBadRequest, result
	}

	result := response.Output{
		Message: "Favorite Found",
		Time:    time.Now(),
		Data:    favoriteData,
	}
	return fiber.StatusOK, result
}
