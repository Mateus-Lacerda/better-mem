package service

import (
	"better-mem/internal/core"
	"better-mem/internal/repository"
	"better-mem/internal/uow"
	"context"
	"time"
)


type MemoryManagementService struct {
	uow uow.UnitOfWork[int]
}

func NewMemoryManagementService(uow uow.UnitOfWork[int]) *MemoryManagementService {
	return &MemoryManagementService{uow: uow}
}

func (s *MemoryManagementService) FindAndDeactivate(
	ctx context.Context,
	chatId string,
	ageLimitHours int,
	minimalRelevance int,
) (int, error) {
	return s.uow.Do(ctx, func(repos repository.AllRepositories) (int, error) {
		endTimeWindow := time.Until(time.Now().Add(time.Duration(ageLimitHours) * time.Hour))
		memories, err := repos.ShortTermMemory.GetElligibleForDeactivation(
			ctx, chatId, endTimeWindow, minimalRelevance,
		)
		deactivated := 0
		if err != nil {
			return deactivated, err
		}
		for _, memory := range memories {
			if err := repos.ShortTermMemory.Deactivate(ctx, chatId, memory.Id); err != nil {
				return deactivated, err
			}
			deactivated++
		}
		return deactivated, nil
	})
}

func (s *MemoryManagementService) FindAndPromote(
	ctx context.Context,
	chatId string,
	minimalRelevance int,
	longTermThreshold float32,
) (int, error) {
	return s.uow.Do(ctx, func(repos repository.AllRepositories) (int, error) {
		memories, err := repos.ShortTermMemory.GetElligibleForPromotion(
			ctx, chatId, minimalRelevance,
		)
		if err != nil {
			return 0, err
		}
		for _, memory := range memories {
			memory, err := repos.LongTermMemory.Create(
				ctx,
				&core.NewLongTermMemory{
					Memory: memory.Memory,
					ChatId: memory.ChatId,
					AccessCount: memory.AccessCount,
					CreatedAt: memory.CreatedAt,
					Active: true,
				},
			)
			if err != nil {
				return 0, err
			}
			if err := repos.ShortTermMemory.Deactivate(ctx, chatId, memory.Id); err != nil {
				return 0, err
			}
		}
		return len(memories), nil
	})
}
