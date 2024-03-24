package mongoDB

import (
	"context"

	"github.com/ez-as/ironlink-base-lib/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBClient struct {
	Client *mongo.Client
	DBName string
}

func NewClient(ctx context.Context, atlasURI, dbName string) *MongoDBClient {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(atlasURI))
	if err != nil {
		logging.Logger().Fatal("error connecting to mongodb", zap.Error(err))
	}

	return &MongoDBClient{
		Client: client,
		DBName: dbName,
	}
}
