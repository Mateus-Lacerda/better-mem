package qdrant

import (
	"better-mem/internal/config"
	"context"
	"log/slog"
	"sync"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	*qdrant.Client
}

var (
	once sync.Once
)

// Get Qdrant client
func NewQdrantClient() *QdrantClient {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   config.Database.QdrantHost,
		Port:   config.Database.QdrantPort,
		UseTLS: false,
	})
	if err != nil {
		slog.Error("Error connecting to Qdrant", "error", err)
		panic(err)
	}
	qdrantClient := &QdrantClient{
		Client: client,
	}
	slog.Info("Qdrant client created")
	return qdrantClient
}

// Check if collection exists
func collectionExists(
	ctx context.Context,
	client *QdrantClient,
) (bool, error) {
	existsRequests := qdrant.CollectionExistsRequest{
		CollectionName: config.Database.DefaultCollectionName,
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
		Size:     config.Database.DefaultVectorSize,
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

// Tests the qdrant collection and ensures the collection exists
func TestQdrant() error {
	ctx := context.Background()
	defer ctx.Done()

	client := NewQdrantClient()
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
