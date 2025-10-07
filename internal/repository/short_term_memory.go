package repository

import (
	"better-mem/internal/core"
	"context"
	"fmt"
	"time"
)

type ShortTermMemoryRepository interface {
	Create(ctx context.Context, memory *core.NewShortTermMemory) (*core.ShortTermMemory, error)
	GetByChatId(
		ctx context.Context,
		chatId string,
		limit int,
		offset int,
	) (*core.ShortTermMemoryArray, error)
	GetById(ctx context.Context, chatId string, memoryId string) (*core.ShortTermMemory, error)
	GetScored(ctx context.Context, chatId string, memoriesIds []string) ([]*core.ScoredMemory, error)
	RegisterUsage(ctx context.Context, chatId string, memoryId string) error
	Deactivate(ctx context.Context, chatId string, memoryId string) error
	Merge(ctx context.Context, chatId string, memoryId string, otherMemory string) (*core.ShortTermMemory, error)
	GetElligibleForDeactivation(ctx context.Context, chatId string, window time.Duration, minimalRelevance int) ([]*core.ShortTermMemory, error)
	GetElligibleForPromotion(ctx context.Context, chatId string, minimalRelevance int) ([]*core.ShortTermMemory, error)
}

type MockShortTermMemoryRepository struct {
	ShortTermMemories []*core.ShortTermMemory
}

func NewMockShortTermMemoryRepository() *MockShortTermMemoryRepository {
	return &MockShortTermMemoryRepository{
		ShortTermMemories: []*core.ShortTermMemory{},
	}
}

func (m *MockShortTermMemoryRepository) Create(ctx context.Context, memory *core.NewShortTermMemory) (*core.ShortTermMemory, error) {
	id := len(m.ShortTermMemories) + 1
	m.ShortTermMemories = append(m.ShortTermMemories, &core.ShortTermMemory{
		Id: fmt.Sprintf("%d", id),
		Memory:      memory.Memory,
		ChatId:      memory.ChatId,
		AccessCount: memory.AccessCount,
		MergeCount:  memory.MergeCount,
		Merged:      memory.Merged,
		CreatedAt:   memory.CreatedAt,
		Active:      memory.Active,
	})
	return m.ShortTermMemories[len(m.ShortTermMemories)-1], nil
}

func (m *MockShortTermMemoryRepository) GetByChatId(
	ctx context.Context,
	chatId string,
	limit int,
	offset int,
) (*core.ShortTermMemoryArray, error) {
	memories := []*core.ShortTermMemory{}
	count := 0
	skip := offset
	for _, memory := range m.ShortTermMemories {
		if skip > 0 {
			skip--
			continue
		}
		if count >= limit {
			break
		}
		if memory.ChatId == chatId {
			memories = append(memories, memory)
		}
	}
	return &core.ShortTermMemoryArray{
		Memories: m.ShortTermMemories,
	}, nil
}

func (m *MockShortTermMemoryRepository) GetById(ctx context.Context, chatId string, memoryId string) (*core.ShortTermMemory, error) {
	for _, memory := range m.ShortTermMemories {
		if memory.ChatId == chatId && memory.Id == memoryId {
			return memory, nil
		}
	}
	return nil, nil
}

func (m *MockShortTermMemoryRepository) GetScored(ctx context.Context, chatId string, memoriesIds []string) ([]*core.ScoredMemory, error) {
	return []*core.ScoredMemory{}, nil
}

func (m *MockShortTermMemoryRepository) RegisterUsage(ctx context.Context, chatId string, memoryId string) error {
	return nil
}

func (m *MockShortTermMemoryRepository) Deactivate(ctx context.Context, chatId string, memoryId string) error {
	return nil
}

func (m *MockShortTermMemoryRepository) Merge(ctx context.Context, chatId string, memoryId string, otherMemoryId string) (*core.ShortTermMemory, error) {
	return nil, nil
}
