package vector

import (
	"better-mem/internal/core"
	"context"
)

type MemoryVectorRepository interface {
	Create(
		ctx context.Context,
		chatId string,
		vectors []float32,
		memoryType core.MemoryTypeEnum,
		memoryId string,
	) error
	Search(
		ctx context.Context,
		chatId string,
		vector []float32,
		limit int,
		threshold float32,
	) (*[]core.ScoredMemoryVector, error)
	Deactivate(ctx context.Context, chatId string, id string) error
	DeactivateAll(ctx context.Context, chatId string) error
}
