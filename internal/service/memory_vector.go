package service

import (
	"context"
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"github.com/Mateus-Lacerda/better-mem/internal/repository/vector"
)


type MemoryVectorService struct {
	repo vector.MemoryVectorRepository
}

func NewMemoryVectorService(repo vector.MemoryVectorRepository) *MemoryVectorService {
	return &MemoryVectorService{repo: repo}
}

func (s *MemoryVectorService) CreateMemoryVector(
	ctx context.Context,
	chatId string,
	vectors []float32,
	memoryType core.MemoryTypeEnum,
	memoryId string,
) error {
	return s.repo.Create(ctx, chatId, vectors, memoryType, memoryId)
}

func (s *MemoryVectorService) SearchMemoryVector(
	ctx context.Context,
	chatId string,
	vector []float32,
	limit int,
	threshold float32,
) (*[]core.ScoredMemoryVector, error) {
	return s.repo.Search(ctx, chatId, vector, limit, threshold)
}
