package repository

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"context"
	"log/slog"

	"github.com/Mateus-Lacerda/better-mem/internal/database/mongo"
	"github.com/Mateus-Lacerda/better-mem/internal/repository"

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
func (r *ChatRepository) Create(ctx context.Context, chat *core.NewChat) error {
	dbChat := mongo.Chat{ExternalID: chat.ExternalId}
	result, err := r.InsertOne(ctx, dbChat)
	if IsMongoDuplicateKeyError(err) {
		return core.ChatExternalIdAlreadyExists
	}
	if err != nil {
		return err
	}
	slog.Info("Chat created", "id", result.InsertedID)
	return nil
}

// GetByExternalID implements repository.ChatRepository.
func (r *ChatRepository) GetByExternalID(ctx context.Context, externalID string) (*string, error) {
	result := r.FindOne(ctx, bson.D{{Key: "externalid", Value: externalID}})
	if result.Err() == mongoDriver.ErrNoDocuments {
		return nil, core.ChatNotFound
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	var chat mongo.Chat
	err := result.Decode(&chat)
	if err != nil {
		return nil, err
	}
	return &chat.ID, nil
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
	var dbChats []*mongo.Chat
	err = result.All(ctx, &dbChats)
	if err != nil {
		slog.Error("Error parsing chats", "error", err)
		return nil, err
	}
	var chats []*core.Chat
	for _, chat := range dbChats {
		chats = append(chats, &core.Chat{
			ExternalId: chat.ExternalID,
			ID:         chat.ID,
		})
	}
	return chats, nil
}

var _ repository.ChatRepository = (*ChatRepository)(nil)
