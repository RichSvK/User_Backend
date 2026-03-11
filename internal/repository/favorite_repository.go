package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"stock_backend/internal/model/domainerr"
	"stock_backend/internal/model/entity"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type FavoriteRepository interface {
	Create(favorite *entity.Favorite, ctx context.Context) error
	GetFavorites(userId string, ctx context.Context) ([]string, error)
	AddFavoriteCache(key string, favorite []string, ctx context.Context) error
	RemoveFavorite(userId string, underwriterCode string, ctx context.Context) error
}

type FavoriteRepositoryImpl struct {
	DB      *sql.DB
	RedisDB *redis.Client
}

func NewFavoriteRepository(db *sql.DB, redisDb *redis.Client) FavoriteRepository {
	return &FavoriteRepositoryImpl{
		DB:      db,
		RedisDB: redisDb,
	}
}

func (repository *FavoriteRepositoryImpl) Create(favorite *entity.Favorite, ctx context.Context) error {
	query := `INSERT INTO favorites (userId, underwriterId) VALUES ($1, $2)`

	if _, err := repository.DB.ExecContext(ctx, query,
		favorite.UserID,
		favorite.UnderwriterID,
	); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return domainerr.ErrFavoritesDuplicate
			}
		}
		return err
	}

	// Invalidate the user's favorites cache so the next read gets fresh data
	cacheKey := fmt.Sprintf("favorites:%s", favorite.UserID)
	if err := repository.RedisDB.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Failed to invalidate cache:", err.Error())
	}
	return nil
}

func (repository *FavoriteRepositoryImpl) GetFavorites(userId string, ctx context.Context) ([]string, error) {
	cacheKey := fmt.Sprintf("favorites:%s", userId)
	cachedData, err := repository.RedisDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var favorite []string
		if err := json.Unmarshal([]byte(cachedData), &favorite); err == nil {
			return favorite, nil
		}
	} else if err != redis.Nil {
		return nil, err
	}

	// If the data is not in cache
	rows, err := repository.DB.QueryContext(ctx, `SELECT underwriterId FROM favorites WHERE userId = $1`, userId)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var favorites []entity.Favorite
	for rows.Next() {
		var fav entity.Favorite
		if err := rows.Scan(&fav.UnderwriterID); err != nil {
			return nil, err
		}
		favorites = append(favorites, fav)
	}

	listFavorite := []string{}
	for _, fav := range favorites {
		listFavorite = append(listFavorite, fav.UnderwriterID)
	}

	_ = repository.AddFavoriteCache(cacheKey, listFavorite, ctx)

	return listFavorite, nil
}

func (repository *FavoriteRepositoryImpl) AddFavoriteCache(key string, favorites []string, ctx context.Context) error {
	if jsonData, err := json.Marshal(favorites); err == nil {
		_ = repository.RedisDB.Set(ctx, key, jsonData, 5*time.Minute).Err()
	}
	return nil
}

func (repository *FavoriteRepositoryImpl) RemoveFavorite(userId string, underwriterCode string, ctx context.Context) error {
	// Remove the cache entry
	cacheKey := fmt.Sprintf("favorites:%s", userId)
	if err := repository.RedisDB.Del(ctx, cacheKey).Err(); err != nil {
		return err
	}

	// Delete in database
	result, err := repository.DB.ExecContext(ctx,
		"DELETE FROM favorites WHERE userId = $1 AND underwriterId = $2",
		userId,
		underwriterCode,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domainerr.ErrFavoritesNotFound
	}

	return nil
}
