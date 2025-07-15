package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(ctx context.Context) (*mongo.Client, error) {
	// Set client options
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://admin:12345678@localhost:27017/test?authSource=admin"))

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
