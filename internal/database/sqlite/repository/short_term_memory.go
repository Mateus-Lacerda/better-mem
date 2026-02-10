package repository

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"github.com/Mateus-Lacerda/better-mem/internal/database/sqlite"
	"github.com/Mateus-Lacerda/better-mem/internal/repository"
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShortTermMemoryRepository struct {
	*gorm.DB
	helper ShortTermMemoryHelper
}

func NewShortTermMemoryRepository() ShortTermMemoryRepository {
	db := sqlite.GetDb()
	return ShortTermMemoryRepository{
		DB:     db,
		helper: ShortTermHelper,
	}
}

func ShortTermMemoryRepositoryWithTransaction(db *gorm.DB) *ShortTermMemoryRepository {
	return &ShortTermMemoryRepository{
		DB:     db,
		helper: ShortTermHelper,
	}
}
func (s *ShortTermMemoryRepository) G() gorm.Interface[sqlite.ShortTermMemory] {
	return gorm.G[sqlite.ShortTermMemory](s.DB)
}

// Create implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) Create(ctx context.Context, memory *core.NewShortTermMemory) (*core.ShortTermMemory, error) {
	dbMemory := s.helper.SchemaToDbModel(memory)
	if err := s.G().Create(ctx, dbMemory); err != nil {
		return nil, err
	}
	createdMemory := &core.ShortTermMemory{
		Id:          dbMemory.ID,
		Memory:      memory.Memory,
		ChatId:      memory.ChatId,
		AccessCount: memory.AccessCount,
		MergeCount:  memory.MergeCount,
		Merged:      memory.Merged,
		CreatedAt:   memory.CreatedAt,
		Active:      memory.Active,
	}
	return createdMemory, nil
}

// Deactivate implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) Deactivate(ctx context.Context, chatId string, memoryId string) error {
	if _, err := s.G().
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		Update(ctx, "active", false); err != nil {
		return err
	}
	return nil
}

// GetByChatId implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) GetByChatId(ctx context.Context, chatId string, limit int, offset int) (*core.ShortTermMemoryArray, error) {
	dbMemories, err := s.G().
		Where("chat_id = ?", chatId).Find(ctx)
	if err != nil {
		return nil, err
	}
	var memories []*core.ShortTermMemory
	total := len(dbMemories)
	for _, m := range dbMemories {
		memories = append(memories, s.helper.DbModelToSchema(&m))
	}
	return &core.ShortTermMemoryArray{
		Memories: memories,
		Total:    int(total),
	}, nil
}

// GetById implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) GetById(ctx context.Context, chatId string, memoryId string) (*core.ShortTermMemory, error) {
	dbMemory, err := s.G().
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return s.helper.DbModelToSchema(&dbMemory), nil
}

// GetScored implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) GetScored(
	ctx context.Context,
	chatId string,
	memoriesIds []string,
) ([]*core.ScoredMemory, error) {
	var memories []*core.ScoredMemory
	if len(memoriesIds) == 0 {
		return memories, nil
	}
	dbMemories, err := s.G().Where("id IN ?", memoriesIds).Find(ctx)
	if err != nil {
		slog.Error("failed to get memories", "error", err)
		return nil, err
	}
	var rawMemories []core.ShortTermMemory
	for _, m := range dbMemories {
		rawMemories = append(rawMemories, *s.helper.DbModelToSchema(&m))
	}

	maxAccessCount, maxMergeCount := s.helper.GetMaxCounts(rawMemories)
	now := time.Now().Unix()
	for _, memory := range rawMemories {
		score, err := s.helper.CalculateScore(
			memory, maxAccessCount, maxMergeCount, now,
		)
		if err != nil {
			return nil, err
		}
		memories = append(memories, &core.ScoredMemory{
			Id:             memory.Id,
			Score:          score,
			Text:           memory.Memory,
			MemoryType:     core.ShortTerm,
			CreatedAt:      memory.CreatedAt,
			RelatedContext: memory.RelatedContext,
		})
	}
	return memories, nil
}

// Merge implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) Merge(
	ctx context.Context,
	chatId string,
	memoryId string,
	otherMemory string,
	otherMemoryRelatedContext []core.MessageRelatedContext,
) (*core.ShortTermMemory, error) {
	// We will just use the newest memory text, and increment the merge count
	// TODO: Store merges in a separate collection for data analysis
	var updatedMemory sqlite.ShortTermMemory
	result := s.DB.Model(&updatedMemory).
		Clauses(clause.Returning{}).
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		Updates(map[string]any{
			"memory":          otherMemory,
			"related_context": otherMemoryRelatedContext,
			"merge_count":     gorm.Expr("merge_count + ?", 1),
		})
	if result.Error != nil {
		return nil, result.Error
	}
	return s.helper.DbModelToSchema(&updatedMemory), nil
}

// RegisterUsage implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) RegisterUsage(ctx context.Context, chatId string, memoryId string) error {
	if _, err := s.G().
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		Update(ctx, "access_count", gorm.Expr("access_count + ?", 1)); err != nil {
		return err
	}
	return nil
}

// GetElligibleForDeactivation implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) GetElligibleForDeactivation(
	ctx context.Context,
	chatId string,
	window time.Duration,
	minimalRelevance int,
) ([]*core.ShortTermMemory, error) {
	dbMemories, err := s.G().
		Where(
			"chat_id = ? AND active AND created_at >= ? AND access_count + merge_count < ?",
			chatId, time.Now().Add(-window), minimalRelevance,
		).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	var memories []*core.ShortTermMemory
	for _, m := range dbMemories {
		memories = append(memories, s.helper.DbModelToSchema(&m))
	}
	return memories, nil
}

// GetElligibleForPromotion implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) GetElligibleForPromotion(
	ctx context.Context, chatId string, minimalRelevance int,
) ([]*core.ShortTermMemory, error) {
	dbMemories, err := s.G().
		Where(
			"chat_id = ? AND active AND access_count + merge_count >= ?",
			chatId, minimalRelevance,
		).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	var memories []*core.ShortTermMemory
	for _, m := range dbMemories {
		memories = append(memories, s.helper.DbModelToSchema(&m))
	}
	return memories, nil
}

// DeactivateAll implements [repository.ShortTermMemoryRepository]
func (s ShortTermMemoryRepository) DeactivateAll(ctx context.Context, chatId string) error {
	if _, err := s.G().
		Where("chat_id = ?", chatId).
		Update(ctx, "active", false); err != nil {
		return err
	}
	return nil
}

var _ repository.ShortTermMemoryRepository = (*ShortTermMemoryRepository)(nil)
