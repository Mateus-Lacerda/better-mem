package repository

import (
	"better-mem/internal/core"
	"context"
	"fmt"
)

type LongTermMemoryRepository interface {
	Create(ctx context.Context, memory *core.NewLongTermMemory) (*core.LongTermMemory, error)
	GetByChatId(
		ctx context.Context,
		chatId string,
		limit int,
		offset int,
	) (*core.LongTermMemoryArray, error)
	GetById(ctx context.Context, chatId string, memoryId string) (*core.LongTermMemory, error)
	GetScored(ctx context.Context, chatId string, memoriesIds []string) ([]*core.ScoredMemory, error)
	RegisterUsage(ctx context.Context, chatId string, memoryId string) error
	Deactivate(ctx context.Context, chatId string, memoryId string) error
}

type MockLongTermMemoryRepository struct {
	LongTermMemories []*core.LongTermMemory
}

func NewMockLongTermMemoryRepository() *MockLongTermMemoryRepository {
	return &MockLongTermMemoryRepository{
		LongTermMemories: []*core.LongTermMemory{},
	}
}

func (m *MockLongTermMemoryRepository) Create(ctx context.Context, memory *core.NewLongTermMemory) (*core.LongTermMemory, error) {
	id := len(m.LongTermMemories) + 1
	m.LongTermMemories = append(m.LongTermMemories, &core.LongTermMemory{
		Id:          fmt.Sprintf("%d", id),
		Memory:      memory.Memory,
		ChatId:      memory.ChatId,
		AccessCount: memory.AccessCount,
		CreatedAt:   memory.CreatedAt,
		Active:      memory.Active,
	})
	return m.LongTermMemories[len(m.LongTermMemories)-1], nil
}

func (m *MockLongTermMemoryRepository) GetByChatId(
	ctx context.Context,
	chatId string,
	limit int,
	offset int,
) (*core.LongTermMemoryArray, error) {
	memories := []*core.LongTermMemory{}
	count := 0
	skip := offset
	for _, memory := range m.LongTermMemories {
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
	return &core.LongTermMemoryArray{
		Memories: memories,
		Total:    len(memories),
	}, nil
}

func (m *MockLongTermMemoryRepository) GetById(ctx context.Context, chatId string, memoryId string) (*core.LongTermMemory, error) {
	for _, memory := range m.LongTermMemories {
		if memory.Id == memoryId {
			return memory, nil
		}
	}
	return nil, nil
}

func (m *MockLongTermMemoryRepository) GetScored(ctx context.Context, chatId string, memoriesIds []string) (*[]core.ScoredMemory, error) {
	return nil, nil
}

func (m *MockLongTermMemoryRepository) RegisterUsage(ctx context.Context, chatId string, memoryId string) error {
	return nil
}

func (m *MockLongTermMemoryRepository) Deactivate(ctx context.Context, chatId string, memoryId string) error {
	return nil
}
