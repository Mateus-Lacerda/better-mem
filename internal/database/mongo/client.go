package mongo

import (
	"github.com/Mateus-Lacerda/better-mem/internal/config"
	"context"
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	*mongo.Client
}

var (
	lock = &sync.Mutex{}
)

func GetMongoClient() *MongoClient {
	options := options.Client().
		ApplyURI(config.Database.MongoUri).
		SetConnectTimeout(
			time.Duration(config.Database.MongoTimeout) * time.Second,
		)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		slog.Error("Error connecting to MongoDB", "error", err)
		panic(err)
	}
	mongoClient := &MongoClient{
		Client: client,
	}
	return mongoClient
}

func GetMongoDatabase() *mongo.Database {
	return GetMongoClient().Database(config.Database.MongoDatabase)
}

// Ensures MongoDB connection works and creates all collections and indexes
func TestMongo() error {
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
