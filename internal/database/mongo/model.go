package mongo

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chat struct {
	ID         string `bson:"_id,omitempty"`
	ExternalID string `bson:"external_id"`
}

type chatConfig struct {
	CollectionName string
}

func ChatConfig() chatConfig {
	return chatConfig{
		CollectionName: "chats",
	}
}

type LongTermMemory struct {
	ID          string `bson:"_id,omitempty"`
	Memory      string `bson:"memory"`
	ChatID      string `bson:"chat_id"`
	AccessCount int    `bson:"access_count"`
	CreatedAt   string `bson:"created_at"`
	Active      bool   `bson:"active"`
}

type ShortTermMemory struct {
	ID          string `bson:"_id,omitempty"`
	Memory      string `bson:"memory"`
	ChatID      string `bson:"chat_id"`
	AccessCount int    `bson:"access_count"`
	MergeCount  int    `bson:"merge_count"`
	Merged      bool   `bson:"merged"`
	CreatedAt   string `bson:"created_at"`
	Active      bool   `bson:"active"`
}

type MemoryConfig struct {
	CollectionName string
	Indexes        mongo.IndexModel
}

func LongTermMemoryConfig() MemoryConfig {
	indexes := mongo.IndexModel{
		Keys: bson.D{
			{Key: "chat_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
	}

	return MemoryConfig{
		CollectionName: "long_term_memory",
		Indexes:        indexes,
	}
}

func ShortTermMemoryConfig() MemoryConfig {
	indexes := mongo.IndexModel{
		Keys: bson.D{
			{Key: "chat_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
	}

	return MemoryConfig{
		CollectionName: "short_term_memory",
		Indexes:        indexes,
	}
}

func CreateCollections(db mongo.Database) error {
	ctx := context.Background()
	defer ctx.Done()

	longTermMemoryConfig := LongTermMemoryConfig()
	shortTermMemoryConfig := ShortTermMemoryConfig()
	chatConfig := ChatConfig()
	collections := []string{
		longTermMemoryConfig.CollectionName,
		shortTermMemoryConfig.CollectionName,
		chatConfig.CollectionName,
	}
	for _, collection := range collections {
		err := db.CreateCollection(ctx, collection)
		if err != nil {
			slog.Error(
				"failed to create collection",
				"collection", collection,
				"error", err,
			)
			return err
		}
	}

	return nil
}

func CreateIndexes(db mongo.Database) error {
	longTermMemoryConfig := LongTermMemoryConfig()
	shortTermMemoryConfig := ShortTermMemoryConfig()

	ctx := context.Background()
	defer ctx.Done()

	_, err := db.Collection(
		longTermMemoryConfig.CollectionName,
	).Indexes().CreateOne(
		ctx, longTermMemoryConfig.Indexes,
	)
	if err != nil {
		slog.Error(
			"failed to create indexes for long term memory",
			"error", err,
		)
		return err
	}

	_, err = db.Collection(
		shortTermMemoryConfig.CollectionName,
	).Indexes().CreateOne(
		ctx, shortTermMemoryConfig.Indexes,
	)
	if err != nil {
		slog.Error(
			"failed to create indexes for short term memory",
			"error", err,
		)
		return err
	}

	return nil
}
