package qdrant

import (
	"context"
	"log/slog"
	"sync"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	*qdrant.Client
}

const (
	qdrantHost     string = "qdrant"
	qdrantPort     int    = 6334
	defaultVectorSize uint64  = 384
	DefaultCollectionName string = "better-mem-default"
)

var (
	once sync.Once
	qdrantClient *QdrantClient
)

// Get Qdrant client
func GetQdrantClient() *QdrantClient {
	// TODO: Use dependency injection instead of singleton
	once.Do(func() {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:     qdrantHost,
			Port:     qdrantPort,
			UseTLS:   false,
		})
		if err != nil {
			slog.Error("Error connecting to Qdrant", "error", err)
			panic(err)
		}
		qdrantClient = &QdrantClient{
			Client: client,
		}
		slog.Info("Qdrant client created")
	})
	return qdrantClient
}

// Check if collection exists
func collectionExists(
	ctx context.Context,
	client *QdrantClient,
) (bool, error) {
	existsRequests := qdrant.CollectionExistsRequest{
		CollectionName: DefaultCollectionName,
	}
	existsResponse, err := client.GetCollectionsClient().CollectionExists(
		ctx,
		&existsRequests,
	)
	return existsResponse.Result.Exists, err
}

// Create collection
func createCollection(
	ctx context.Context,
	client *QdrantClient,
) error {
	vectorsConfig := qdrant.NewVectorsConfig(&qdrant.VectorParams{
		Size:     defaultVectorSize,
		Distance: qdrant.Distance_Cosine,
	})
	createRequests := qdrant.CreateCollection{
		CollectionName: "better-mem-default",
		VectorsConfig:  vectorsConfig,
	}
	_, err := client.GetCollectionsClient().Create(
		ctx,
		&createRequests,
	)
	return err
}

// Ensure collection exists
func EnsureCollection(
	ctx context.Context,
	client *QdrantClient,
) error {
	exists, err := collectionExists(ctx, client)
	if err != nil {
		return err
	}

	switch exists {
	case true:
		return nil
	case false:
		return createCollection(ctx, client)
	}
	return nil
}

// Creates a global Qdrant client and ensures the collection exists
func SetupQdrant() error {
	ctx := context.Background()
	defer ctx.Done()

	client := GetQdrantClient()
	slog.Info("Connecting to Qdrant")
	_, err := client.HealthCheck(ctx)
	if err != nil {
		slog.Error("Error pinging Qdrant", "error", err)
		return err
	}
	slog.Info("Qdrant connected")
	err = EnsureCollection(ctx, client)
	if err != nil {
		slog.Error("Error ensuring collection", "error", err)
		return err
	}

	slog.Info("Qdrant setup complete")
	return nil
}
