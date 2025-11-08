package repository

import (
	"better-mem/internal/core"
	"better-mem/internal/database/mongo"
	"better-mem/internal/repository"
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortTermMemoryRepository struct {
	*mongoDriver.Collection
	helper ShortTermMemoryHelper
}

func NewShortTermMemoryRepository() ShortTermMemoryRepository {
	collectionName := mongo.ShortTermMemoryConfig().CollectionName
	database := mongo.GetMongoDatabase()
	return ShortTermMemoryRepository{
		database.Collection(collectionName),
		ShortTermHelper,
	}
}

// Create implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) Create(ctx context.Context, memory *core.NewShortTermMemory) (*core.ShortTermMemory, error) {
	res, err := s.InsertOne(ctx, memory)
	if err != nil {
		return nil, err
	}
	createdMemory := &core.ShortTermMemory{
		Id:          res.InsertedID.(primitive.ObjectID).Hex(),
		Memory:      memory.Memory,
		ChatId:      memory.ChatId,
		AccessCount: memory.AccessCount,
		MergeCount:  memory.MergeCount,
		Merged:      memory.Merged,
		CreatedAt:   memory.CreatedAt,
		Active:      memory.Active,
	}
	return createdMemory, nil
}

// Deactivate implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) Deactivate(ctx context.Context, chatId string, memoryId string) error {
	memoryIdObjectId, err := primitive.ObjectIDFromHex(memoryId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": memoryIdObjectId, "chatid": chatId}
	update := bson.M{"$set": bson.M{"active": false}}
	_, err = s.UpdateOne(ctx, filter, update)
	return err
}

// GetByChatId implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) GetByChatId(ctx context.Context, chatId string, limit int, offset int) (*core.ShortTermMemoryArray, error) {
	filter := bson.M{"chatid": chatId, "active": true}
	cursor, err := s.Find(
		ctx,
		filter,
		options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	var memories []*core.ShortTermMemory
	total, err := s.CountDocuments(ctx, filter)
	err = cursor.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	return &core.ShortTermMemoryArray{
		Memories: memories,
		Total:    int(total),
	}, nil
}

// GetById implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) GetById(ctx context.Context, chatId string, memoryId string) (*core.ShortTermMemory, error) {
	memoryIdObjectId, err := primitive.ObjectIDFromHex(memoryId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": memoryIdObjectId, "chatid": chatId}
	result := s.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var memory core.ShortTermMemory
	err = result.Decode(&memory)
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// GetScored implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) GetScored(
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
		"_id":    bson.M{"$in": objectIds},
		"active": true,
	}
	cursor, err := s.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rawMemories []core.ShortTermMemoryModel
	err = cursor.All(ctx, &rawMemories)
	if err != nil {
		return nil, err
	}

	maxAccessCount, maxMergeCount := s.helper.GetMaxCounts(rawMemories)
	now := time.Now().Unix()
	var memories []*core.ScoredMemory

	for _, memory := range rawMemories {
		score, err := s.helper.CalculateScore(
			memory, maxAccessCount, maxMergeCount, now,
		)
		if err != nil {
			return nil, err
		}
		memories = append(memories, &core.ScoredMemory{
			Score:      score,
			Text:       memory.Memory,
			MemoryType: core.ShortTerm,
			CreatedAt:  memory.CreatedAt,
			RelatedContext: memory.RelatedContext,
		})
	}
	return memories, nil
}

// Merge implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) Merge(
	ctx context.Context,
	chatId string,
	memoryId string,
	otherMemory string,
	otherMemoryRelatedContext []core.MessageRelatedContext,
) (*core.ShortTermMemory, error) {
	// We will just use the newest memory text, and increment the merge count
	// TODO: Store merges in a separate collection for data analysis
	oldMemory, err := s.GetById(ctx, chatId, memoryId)
	if err != nil {
		return nil, err
	}
	oldMemory.Memory = otherMemory
	oldMemory.RelatedContext = otherMemoryRelatedContext
	oldMemory.MergeCount++
	memoryIdObjectId, err := primitive.ObjectIDFromHex(memoryId)
	filter := bson.M{"_id": memoryIdObjectId}
	update := bson.M{"$set": oldMemory}
	_, err = s.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return oldMemory, nil
}

// RegisterUsage implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) RegisterUsage(ctx context.Context, chatId string, memoryId string) error {
	memory, err := s.GetById(ctx, chatId, memoryId)
	if err != nil {
		return err
	}
	memory.AccessCount++
	slog.Info("Memory accessed", "memory", memory)
	memoryIdObjectId, err := primitive.ObjectIDFromHex(memory.Id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": memoryIdObjectId}
	update := bson.M{"$set": bson.M{"accesscount": memory.AccessCount}}
	_, err = s.UpdateOne(ctx, filter, update)
	return err
}

// GetElligibleForDeactivation implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) GetElligibleForDeactivation(
	ctx context.Context,
	chatId string,
	window time.Duration,
	minimalRelevance int,
) ([]*core.ShortTermMemory, error) {
	filter := bson.M{
		"chatid": chatId,
		"active": true,
		"createdat": bson.M{
			"$lt": time.Now().Add(-window),
		},
		"$expr": bson.M{
			"$lt": []any{
				bson.M{"$add": []any{"$accesscount", "$mergecount"}},
				minimalRelevance,
			},
		},
	}
	cursor, err := s.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var memories []*core.ShortTermMemory
	err = cursor.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	return memories, nil
}

// GetElligibleForPromotion implements repository.ShortTermMemoryRepository.
func (s ShortTermMemoryRepository) GetElligibleForPromotion(
	ctx context.Context, chatId string, minimalRelevance int,
) ([]*core.ShortTermMemory, error) {
	filter := bson.M{
		"chatid": chatId,
		"active": true,
		"$expr": bson.M{
			"$gte": []any{
				bson.M{"$add": []any{"$accesscount", "$mergecount"}},
				minimalRelevance,
			},
		},
	}
	cursor, err := s.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var memories []*core.ShortTermMemory
	err = cursor.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	return memories, nil
}

var _ repository.ShortTermMemoryRepository = (*ShortTermMemoryRepository)(nil)
