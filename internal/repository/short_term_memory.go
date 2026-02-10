package repository

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"context"
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
	Merge(ctx context.Context, chatId string, memoryId string, otherMemory string, otherMemoryRelatedContext []core.MessageRelatedContext) (*core.ShortTermMemory, error)
	GetElligibleForDeactivation(ctx context.Context, chatId string, window time.Duration, minimalRelevance int) ([]*core.ShortTermMemory, error)
	GetElligibleForPromotion(ctx context.Context, chatId string, minimalRelevance int) ([]*core.ShortTermMemory, error)
	DeactivateAll(ctx context.Context, chatId string) error
}
