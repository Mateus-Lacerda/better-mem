package repository

import (
	"context"
	"log/slog"
	"time"
	"better-mem/internal/core"
	"better-mem/internal/database/mongo"
	"better-mem/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LongTermMemoryRepository struct {
	*mongoDriver.Collection
	helper LongTermMemoryHelper
}

func NewLongTermMemoryRepository() *LongTermMemoryRepository {
	collectionName := mongo.LongTermMemoryConfig().CollectionName
	database := mongo.GetMongoDatabase()
	return &LongTermMemoryRepository{
		Collection: database.Collection(collectionName),
		helper:     LongTermHelper,
	}
}

// Create implements repository.LongTermMemoryRepository.
func (l *LongTermMemoryRepository) Create(ctx context.Context, memory *core.NewLongTermMemory) (*core.LongTermMemory, error) {
	res, err := l.InsertOne(ctx, memory)
	if err != nil {
		return nil, err
	}
	createdMemory := &core.LongTermMemory{
		Id:                res.InsertedID.(primitive.ObjectID).Hex(),
		NewLongTermMemory: *memory,
	}
	return createdMemory, nil
}

// Deactivate implements repository.LongTermMemoryRepository.
func (l *LongTermMemoryRepository) Deactivate(ctx context.Context, chatId string, memoryId string) error {
	filter := bson.M{"chatid": chatId, "_id": memoryId}
	update := bson.M{"$set": bson.M{"active": false}}
	_, err := l.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetByChatId implements repository.LongTermMemoryRepository.
func (l *LongTermMemoryRepository) GetByChatId(ctx context.Context, chatId string, limit int, offset int) (*core.LongTermMemoryArray, error) {
	filter := bson.M{"chatid": chatId, "active": true}
	cursor, err := l.Find(
		ctx,
		filter,
		options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	var memories []*core.LongTermMemory
	total, err := l.CountDocuments(ctx, filter)
	err = cursor.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	return &core.LongTermMemoryArray{
		Memories: memories,
		Total:    int(total),
	}, nil
}

// GetById implements repository.LongTermMemoryRepository.
func (l *LongTermMemoryRepository) GetById(ctx context.Context, chatId string, memoryId string) (*core.LongTermMemory, error) {
	filter := bson.M{"_id": memoryId, "chat_id": chatId}
	result := l.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var memory core.LongTermMemory
	err := result.Decode(&memory)
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// GetScored implements repository.LongTermMemoryRepository.
func (l *LongTermMemoryRepository) GetScored(
	ctx context.Context,
	chatId string,
	memoriesIds []string,
) ([]*core.ScoredMemory, error) {
	if len(memoriesIds) == 0 {
		return []*core.ScoredMemory{}, nil
	}
	var objectIds []primitive.ObjectID
	for _, id := range memoriesIds {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIds = append(objectIds, objectId)
	}
	filter := bson.M{
		"chatid": chatId,
		"_id":     bson.M{"$in": objectIds},
		"active":  true,
	}
	cursor, err := l.Find(ctx, filter)
	if err != nil {
		slog.Error("failed to get memories", "error", err)
		return nil, err
	}
	var rawMemories []core.LongTermMemoryModel
	err = cursor.All(ctx, &rawMemories)
	if err != nil {
		slog.Error("failed to get memories", "error", err)
		return nil, err
	}

	var scoredMemories []*core.ScoredMemory
	maxAge, maxAccessCount, err := l.helper.GetMaxCounts(rawMemories)
	if err != nil {
		slog.Error("failed to get max counts", "error", err)
		return nil, err
	}
	slog.Info("max age", "maxAge", maxAge, "maxAccessCount", maxAccessCount)
	now := time.Now().Unix()
	for _, memory := range rawMemories {
		score, err := l.helper.CalculateScore(
			memory,
			maxAge,
			maxAccessCount,
			now,
		)
		if err != nil {
			return nil, err
		}
		scoredMemories = append(
			scoredMemories, &core.ScoredMemory{
				Text:  memory.Memory,
				Score: score,
				MemoryType: core.LongTerm,
			},
		)
	}

	return scoredMemories, nil
}

// RegisterUsage implements repository.LongTermMemoryRepository.
func (l *LongTermMemoryRepository) RegisterUsage(ctx context.Context, chatId string, memoryId string) error {
	memory, err := l.GetById(ctx, chatId, memoryId)
	if err != nil {
		return err
	}
	memory.AccessCount++
	filter := bson.M{"_id": memory.Id}
	update := bson.M{"$set": bson.M{"access_count": memory.AccessCount}}
	_, err = l.UpdateOne(ctx, filter, update)
	return err
}

var _ repository.LongTermMemoryRepository = (*LongTermMemoryRepository)(nil)
