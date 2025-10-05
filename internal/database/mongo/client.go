package mongo

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"better-mem/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	*mongo.Client
}

var (
	lock        = &sync.Mutex{}
	mongoClient *MongoClient
)

func GetMongoClient() *MongoClient {
	if mongoClient == nil {
		lock.Lock()
		defer lock.Unlock()
		options := options.Client().
			ApplyURI(config.DatabaseConfig.MongoUri).
			SetConnectTimeout(
				time.Duration(config.DatabaseConfig.MongoTimeout) * time.Second,
			)
		ctx := context.Background()
		client, err := mongo.Connect(ctx, options)
		if err != nil {
			slog.Error("Error connecting to MongoDB", "error", err)
			panic(err)
		}
		if mongoClient == nil {
			mongoClient = &MongoClient{
				Client: client,
			}
		}
	}
	return mongoClient
}

func GetMongoDatabase() *mongo.Database {
	return GetMongoClient().Database(config.DatabaseConfig.MongoDatabase)
}

func SetupMongo() error {
	ctx := context.Background()
	defer ctx.Done()

	client := GetMongoClient()
	slog.Info("Connecting to MongoDB")
	err := client.Ping(ctx, nil)
	if err != nil {
		slog.Error("Error pinging MongoDB", "error", err)
		return err
	}
	slog.Info("MongoDB connected")

	database := GetMongoDatabase()

	CreateCollections(*database)
	slog.Info("MongoDB collections created")

	CreateIndexes(*database)
	slog.Info("MongoDB indexes created")

	slog.Info("MongoDB setup complete")
	return nil
}
