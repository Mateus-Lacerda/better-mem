package repository

import (
	"context"
	"log/slog"
	"better-mem/internal/core"
	"better-mem/internal/database/qdrant"
	"better-mem/internal/repository/vector"

	"github.com/google/uuid"
	qdrantClient "github.com/qdrant/go-client/qdrant"
)

type MemoryRepository struct {
	*qdrant.QdrantClient
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		QdrantClient: qdrant.GetQdrantClient(),
	}
}

// Create implements vector.MemoryVectorRepository.
func (m *MemoryRepository) Create(
	ctx context.Context,
	chatId string,
	vectors []float32,
	memoryType core.MemoryTypeEnum,
	memoryId string,
) error {
	payload := core.MemoryPayload{
		ChatId:     chatId,
		MemoryType: memoryType,
		Active:     true,
		MemoryId:   memoryId,
	}
	payloadMap := payload.ToMap()

	payloadValueMap, err := qdrantClient.TryValueMap(payloadMap)
	if err != nil {
		slog.Error("error converting payload to value map", "error", err)
		return err
	}
	qdrantUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	points := []*qdrantClient.PointStruct{{
		Id:      qdrantClient.NewID(qdrantUuid.String()),
		Vectors: qdrantClient.NewVectors(vectors...),
		Payload: payloadValueMap,
	}}
	request := qdrantClient.UpsertPoints{
		CollectionName: qdrant.DefaultCollectionName,
		Points:         points,
	}
	info, err := m.Client.Upsert(ctx, &request)
	slog.Info("upsert", "info", info)
	return err
}

// Search implements vector.MemoryVectorRepository.
func (m *MemoryRepository) Search(
	ctx context.Context,
	chatId string,
	vector []float32,
	limit int,
	threshold float32,
) (*[]core.ScoredMemoryVector, error) {
	filter := qdrantClient.Filter{
		Must: []*qdrantClient.Condition{
			qdrantClient.NewMatchBool("active", true),
			qdrantClient.NewMatchText("chat_id", chatId),
		},
	}
	query := qdrantClient.NewQuery(vector...)
	request := qdrantClient.QueryPoints{
		CollectionName: qdrant.DefaultCollectionName,
		Query:          query,
		Filter:         &filter,
		ScoreThreshold: &threshold,
		WithVectors:    qdrantClient.NewWithVectors(true),
		WithPayload:    qdrantClient.NewWithPayload(true),
	}
	result, err := m.Client.Query(
		ctx, &request,
	)
	if err != nil {
		return nil, err
	}
	var memories []core.ScoredMemoryVector
	for _, point := range result {
		vectorsOutput := point.GetVectors()
		if vectorsOutput == nil {
			continue
		}
		vectors := vectorsOutput.GetVector()

		if vectors == nil || len(vectors.Data) == 0 {
			continue
		}

		score := point.GetScore()
		memoryType := core.MemoryTypeEnum(point.GetPayload()["memory_type"].GetIntegerValue())
		memoryId := point.GetPayload()["memory_id"].GetStringValue()
		active := point.GetPayload()["active"].GetBoolValue()

		memories = append(memories, core.ScoredMemoryVector{
			Id:      point.GetId().String(),
			Vectors: vectors.Data,
			Score:   score,
			Payload: core.MemoryPayload{
				ChatId:     chatId,
				MemoryType: memoryType,
				MemoryId:   memoryId,
				Active:     active,
			},
		})
	}
	return &memories, nil

}

// Deactivate implements vector.MemoryVectorRepository.
func (m *MemoryRepository) Deactivate(ctx context.Context, chatId string, id string) error {
	filter := qdrantClient.Filter{
		Must: []*qdrantClient.Condition{
			qdrantClient.NewMatchText("memory_id", id),
			qdrantClient.NewMatchBool("active", true),
			qdrantClient.NewMatchText("chat_id", chatId),
		},
	}
	request := qdrantClient.SetPayloadPoints{
		CollectionName: qdrant.DefaultCollectionName,
		Payload: map[string]*qdrantClient.Value{
			"active": qdrantClient.NewValueBool(false),
		},
		PointsSelector: qdrantClient.NewPointsSelectorFilter(&filter),
	}
	_, err := m.Client.SetPayload(ctx, &request)
	return err
}

var _ vector.MemoryVectorRepository = (*MemoryRepository)(nil)
