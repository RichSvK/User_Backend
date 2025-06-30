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
	AddFavoriteCache(ctx context.Context, key string, favorite *entity.Favorite) error
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

func (r *FavoriteRepositoryImpl) Create(favorite *entity.Favorite, ctx context.Context) error {
	collection := r.DB.Database("test").Collection("favorites")
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

			r.AddFavoriteCache(ctx, cacheKey, favorite)
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
		r.AddFavoriteCache(ctx, cacheKey, &updatedFavorite)
	}
	return err
}

func (r *FavoriteRepositoryImpl) GetFavorites(userId string, ctx context.Context) (*entity.Favorite, error) {
	cacheKey := fmt.Sprintf("favorites:%s", userId)
	cachedData, err := r.RedisDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var favorite entity.Favorite
		if err := json.Unmarshal([]byte(cachedData), &favorite); err == nil {
			return &favorite, nil
		}
	} else if err != redis.Nil {
		return nil, err
	}

	// If the data is not in cache
	collection := r.DB.Database("test").Collection("favorites")

	var favorite entity.Favorite
	err = collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&favorite)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	r.AddFavoriteCache(ctx, cacheKey, &favorite)

	return &favorite, nil
}

func (r *FavoriteRepositoryImpl) AddFavoriteCache(ctx context.Context, key string, favorite *entity.Favorite) error {
	if jsonData, err := json.Marshal(favorite); err == nil {
		_ = r.RedisDB.Set(ctx, key, jsonData, 1*time.Minute).Err()
	}
	return nil
}
