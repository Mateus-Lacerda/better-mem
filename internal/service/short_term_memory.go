package service

import (
	"better-mem/internal/core"
	"better-mem/internal/repository"
	"context"
	"log/slog"
	"time"
)

type ShortTermMemoryService struct {
	repo     repository.ShortTermMemoryRepository
	chatRepo repository.ChatRepository
}

func NewShortTermMemoryService(
	repo repository.ShortTermMemoryRepository,
	chatRepo repository.ChatRepository,
) *ShortTermMemoryService {
	return &ShortTermMemoryService{repo: repo, chatRepo: chatRepo}
}

func (s *ShortTermMemoryService) Create(
	ctx context.Context,
	text string,
	chatId string,
	relatedContext []core.MessageRelatedContext,
) (*core.ShortTermMemory, error) {
	memory := &core.NewShortTermMemory{
		Memory:         text,
		ChatId:         chatId,
		AccessCount:    0,
		MergeCount:     0,
		Merged:         false,
		CreatedAt:      time.Now(),
		Active:         true,
		RelatedContext: relatedContext,
	}
	return s.repo.Create(ctx, memory)
}

func (s *ShortTermMemoryService) GetByChatId(
	ctx context.Context,
	chatExternalId string,
	limit int,
	offset int,
) (*core.ShortTermMemoryArray, error) {
	chatId, err := s.chatRepo.GetByExternalID(ctx, chatExternalId)
	if err != nil {
		slog.Error("error getting chat id", "error", err)
		return nil, err
	}
	if chatId == nil {
		return nil, core.ChatNotFound
	}
	return s.repo.GetByChatId(ctx, *chatId, limit, offset)
}

func (s *ShortTermMemoryService) GetById(
	ctx context.Context, chatId string, memoryId string,
) (*core.ShortTermMemory, error) {
	return s.repo.GetById(ctx, chatId, memoryId)
}

func (s *ShortTermMemoryService) GetScored(
	ctx context.Context, chatId string, memoriesIds []string,
) ([]*core.ScoredMemory, error) {
	return s.repo.GetScored(ctx, chatId, memoriesIds)
}

func (s *ShortTermMemoryService) RegisterUsage(
	ctx context.Context, chatId string, memoryId string,
) error {
	return s.repo.RegisterUsage(ctx, chatId, memoryId)
}

func (s *ShortTermMemoryService) Merge(
	ctx context.Context,
	chatId string,
	memoryId string,
	otherMemory string,
	otherMemoryRelatedContext []core.MessageRelatedContext,
) (*core.ShortTermMemory, error) {
	return s.repo.Merge(ctx, chatId, memoryId, otherMemory, otherMemoryRelatedContext)
}

func (s *ShortTermMemoryService) Deactivate(
	ctx context.Context, chatId string, memoryId string,
) error {
	return s.repo.Deactivate(ctx, chatId, memoryId)
}
