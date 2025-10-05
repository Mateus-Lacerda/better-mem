package service

import (
	// "context"
	// "better-mem/internal/core"
	"context"
	"log/slog"
	"sort"
	"better-mem/internal/core"
	protos "better-mem/internal/grpc_client"
	"better-mem/internal/repository"
	"better-mem/internal/repository/vector"
)

type MemoryService struct {
	shortTermRepo repository.ShortTermMemoryRepository
	longTermRepo  repository.LongTermMemoryRepository
	vectorRepo    vector.MemoryVectorRepository
}

func NewMemoryService(
	shortTermRepo repository.ShortTermMemoryRepository,
	longTermRepo repository.LongTermMemoryRepository,
	vectorRepo vector.MemoryVectorRepository,
) *MemoryService {
	return &MemoryService{
		shortTermRepo: shortTermRepo,
		longTermRepo:  longTermRepo,
		vectorRepo:    vectorRepo,
	}
}


func (s *MemoryService) Fetch(
	ctx context.Context,
	chatId string,
	query string,
	limit int,
	vectorSearchLimit int,
	threshold float32,
	longTermThreshold float32,
) ([]*core.ScoredMemory, error) {
	slog.Warn("TODO: Add similarity score relevancy to count")
	var memories []*core.ScoredMemory
	vectorService := NewMemoryVectorService(s.vectorRepo)
	embeddings, err := protos.Embed(query)
	if err != nil {
		return memories, err
	}
	similarMemories, err := vectorService.SearchMemoryVector(
		ctx,
		chatId,
		embeddings,
		vectorSearchLimit,
		threshold,
	)
	if err != nil {
		slog.Error("Error searching memory vector", "error", err)
		return memories, err
	}
	if similarMemories == nil {
		slog.Info("No similar memories found")
		return memories, nil
	}
	var shortTermMemories []string
	var longTermMemories []string
	for _, memory := range *similarMemories {
		switch memory.Payload.MemoryType {
		case core.ShortTerm:
			s.shortTermRepo.RegisterUsage(ctx, chatId, memory.Payload.MemoryId)
			shortTermMemories = append(
				shortTermMemories,
				memory.Payload.MemoryId,
			)
		case core.LongTerm:
			if memory.Score >= longTermThreshold {
				s.longTermRepo.RegisterUsage(ctx, chatId, memory.Payload.MemoryId)
				longTermMemories = append(
					longTermMemories,
					memory.Payload.MemoryId,
				)
			}
		}
	}
	scoredSTMemories, err := s.shortTermRepo.GetScored(
		ctx, chatId, shortTermMemories,
	)
	if err != nil {
		slog.Error("Error getting short term memories", "error", err)
		return memories, err
	}
	scoredLTMemories, err := s.longTermRepo.GetScored(
		ctx, chatId, longTermMemories,
	)
	if err != nil {
		slog.Error("Error getting long term memories", "error", err)
		return memories, err
	}
	// Sort by score
	memories = append(scoredSTMemories, scoredLTMemories...)
	if len(memories) == 0 {
		return memories, nil
	}
	if limit > len(memories) {
		limit = len(memories)
	}
	sort.SliceStable(
		memories[:],
		func(i, j int) bool {
			return memories[i].Score > memories[j].Score
		},
	)
	return memories[:limit], nil
}
