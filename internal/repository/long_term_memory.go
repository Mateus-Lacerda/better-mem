package repository

import (
	"better-mem/internal/core"
	"context"
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
	DeactivateAll(ctx context.Context, chatId string) error
}
