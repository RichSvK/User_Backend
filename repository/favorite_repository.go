package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"stock_backend/model/entity"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FavoriteRepository interface {
	Create(favorite *entity.Favorite, ctx context.Context) error
	GetFavorites(userId string, ctx context.Context) (*entity.Favorite, error)
	AddFavoriteCache(key string, favorite *entity.Favorite, ctx context.Context) error
	RemoveFavorite(userId string, ctx context.Context) error
}

type FavoriteRepositoryImpl struct {
	DB      *mongo.Client
	RedisDB *redis.Client
}

func NewFavoriteRepository(db *mongo.Client, redisDb *redis.Client) FavoriteRepository {
	return &FavoriteRepositoryImpl{
		DB:      db,
		RedisDB: redisDb,
	}
}

func (repository *FavoriteRepositoryImpl) Create(favorite *entity.Favorite, ctx context.Context) error {
	collection := repository.DB.Database("test").Collection("favorites")
	collection.FindOne(ctx, bson.M{"userId": favorite.UserID})

	filter := bson.M{"userId": favorite.UserID}

	// Check if the user favorite already exists
	err := collection.FindOne(ctx, filter).Err()
	cacheKey := fmt.Sprintf("favorites:%s", favorite.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Insert New Document
			if _, err := collection.InsertOne(ctx, favorite); err != nil {
				return err
			}

			repository.AddFavoriteCache(cacheKey, favorite, ctx)
			return nil
		}

		return err
	}

	// Update document if user favorite found
	update := bson.M{
		"$addToSet": bson.M{"underwriters": bson.M{"$each": favorite.Underwriters}},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	if _, err = collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	var updatedFavorite entity.Favorite
	err = collection.FindOne(ctx, filter).Decode(&updatedFavorite)
	if err == nil {
		repository.AddFavoriteCache(cacheKey, &updatedFavorite, ctx)
	}
	return err
}

func (repository *FavoriteRepositoryImpl) GetFavorites(userId string, ctx context.Context) (*entity.Favorite, error) {
	cacheKey := fmt.Sprintf("favorites:%s", userId)
	cachedData, err := repository.RedisDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var favorite entity.Favorite
		if err := json.Unmarshal([]byte(cachedData), &favorite); err == nil {
			return &favorite, nil
		}
	} else if err != redis.Nil {
		return nil, err
	}

	// If the data is not in cache
	collection := repository.DB.Database("test").Collection("favorites")

	var favorite entity.Favorite
	err = collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&favorite)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	repository.AddFavoriteCache(cacheKey, &favorite, ctx)

	return &favorite, nil
}

func (repository *FavoriteRepositoryImpl) AddFavoriteCache(key string, favorite *entity.Favorite, ctx context.Context) error {
	if jsonData, err := json.Marshal(favorite); err == nil {
		_ = repository.RedisDB.Set(ctx, key, jsonData, 1*time.Minute).Err()
	}
	return nil
}

func (repository *FavoriteRepositoryImpl) RemoveFavorite(userId string, ctx context.Context) error {
	collection := repository.DB.Database("test").Collection("favorites")
	filter := bson.M{"userId": userId}

	// Remove the cache entry
	cacheKey := fmt.Sprintf("favorites:%s", userId)
	if err := repository.RedisDB.Del(ctx, cacheKey).Err(); err != nil {
		return err
	}

	// Remove the document from MongoDB
	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
