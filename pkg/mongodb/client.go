package mongoDB

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	Client *mongo.Client
	DBName string
}

func NewClient(ctx context.Context, atlasURI, dbName string) (*MongoDBClient, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(atlasURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}

	return &MongoDBClient{
		Client: client,
		DBName: dbName,
	}, nil
}
