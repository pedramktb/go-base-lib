package mongoDB

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient is a struct that contains a MongoDB client and the name of the database
type MongoDBClient struct {
	Client *mongo.Client
	DBName string
}

// NewClient creates a new MongoDB client with the given URI and database name
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
