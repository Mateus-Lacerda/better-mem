package service

import (
	"context"
	"time"
	"better-mem/internal/core"
	"better-mem/internal/repository"
)

type LongTermMemoryService struct {
	repo repository.LongTermMemoryRepository
}

func NewLongTermMemoryService(
	repo repository.LongTermMemoryRepository,
) *LongTermMemoryService {
	return &LongTermMemoryService{repo: repo}
}

func (s *LongTermMemoryService) Create(
	ctx context.Context,
	text,
	chatId string,
	relatedContext []core.MessageRelatedContext,
) (*core.LongTermMemory, error) {
	memory := &core.NewLongTermMemory{
		Memory: text,
		ChatId: chatId,
		AccessCount: 0,
		CreatedAt: time.Now(),
		Active: true,
		RelatedContext: relatedContext,
	}
	return s.repo.Create(ctx, memory)
}

func (s *LongTermMemoryService) GetByChatId(
	ctx context.Context,
	chatId string,
	limit int,
	offset int,
) (*core.LongTermMemoryArray, error) {
	return s.repo.GetByChatId(ctx, chatId, limit, offset)
}

func (s *LongTermMemoryService) GetById(
	ctx context.Context, chatId string, memoryId string,
) (*core.LongTermMemory, error) {
	return s.repo.GetById(ctx, chatId, memoryId)
}

func (s *LongTermMemoryService) GetScored(
	ctx context.Context, chatId string, memoriesIds []string,
) ([]*core.ScoredMemory, error) {
	return s.repo.GetScored(ctx, chatId, memoriesIds)
}

func (s *LongTermMemoryService) RegisterUsage(
	ctx context.Context, chatId string, memoryId string,
) error {
	return s.repo.RegisterUsage(ctx, chatId, memoryId)
}

func (s *LongTermMemoryService) Deactivate(
	ctx context.Context, chatId string, memoryId string,
) error {
	return s.repo.Deactivate(ctx, chatId, memoryId)
}
