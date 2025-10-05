package repository

import (
	"context"
	"log/slog"
	"better-mem/internal/core"

	"better-mem/internal/database/mongo"
	"better-mem/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository struct {
	*mongoDriver.Collection
}

func NewChatRepository() *ChatRepository {
	collectionName := mongo.ChatConfig().CollectionName
	database := mongo.GetMongoDatabase()
	return &ChatRepository{
		Collection: database.Collection(collectionName),
	}
}

// Create implements repository.ChatRepository.
func (r *ChatRepository) Create(ctx context.Context, chat *core.Chat) error {
	result, err := r.InsertOne(ctx, chat)
	if err != nil {
		return err
	}
	slog.Info("Chat created", "id", result.InsertedID)
	return nil
}

// GetAll implements repository.ChatRepository.
func (r *ChatRepository) GetAll(ctx context.Context) ([]*core.Chat, error) {
	result, err := r.Find(
		ctx, bson.D{},
	)
	if err != nil {
		slog.Error("Error getting all chats", "error", err)
		return nil, err
	}
	var chats []*core.Chat
	err = result.All(ctx, &chats)
	if err != nil {
		slog.Error("Error parsing chats", "error", err)
		return nil, err
	}
	return chats, nil
}

var _ repository.ChatRepository = (*ChatRepository)(nil)
