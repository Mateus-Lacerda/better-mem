package repository

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	sqlite "github.com/Mateus-Lacerda/better-mem/internal/database/sqlite"
	"github.com/Mateus-Lacerda/better-mem/internal/repository"
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type LongTermMemoryRepository struct {
	*gorm.DB
	helper LongTermMemoryHelper
}

func NewLongTermMemoryRepository() *LongTermMemoryRepository {
	db := sqlite.GetDb()
	return &LongTermMemoryRepository{
		DB:     db,
		helper: LongTermHelper,
	}
}

func LongTermMemoryRepositoryWithTransaction(db *gorm.DB) *LongTermMemoryRepository {
	return &LongTermMemoryRepository{
		DB:     db,
		helper: LongTermHelper,
	}
}

func (l *LongTermMemoryRepository) G() gorm.Interface[sqlite.LongTermMemory] {
	return gorm.G[sqlite.LongTermMemory](l.DB)
}

// Create implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) Create(ctx context.Context, memory *core.NewLongTermMemory) (*core.LongTermMemory, error) {
	dbMemory := l.helper.SchemaToDbModel(memory)
	if err := l.G().Create(ctx, dbMemory); err != nil {
		return nil, err
	}

	createdMemory := &core.LongTermMemory{
		Id:          dbMemory.ID,
		Memory:      memory.Memory,
		ChatId:      memory.ChatId,
		AccessCount: memory.AccessCount,
		CreatedAt:   memory.CreatedAt,
		Active:      memory.Active,
	}
	return createdMemory, nil
}

// Deactivate implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) Deactivate(ctx context.Context, chatId string, memoryId string) error {
	if _, err := l.G().
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		Update(ctx, "active", false); err != nil {
		return err
	}
	return nil
}

// GetByChatId implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) GetByChatId(ctx context.Context, chatId string, limit int, offset int) (*core.LongTermMemoryArray, error) {
	dbMemories, err := l.G().
		Where("chat_id = ?", chatId).Find(ctx)
	if err != nil {
		return nil, err
	}
	var memories []*core.LongTermMemory
	total := len(dbMemories)
	for _, m := range dbMemories {
		memories = append(memories, l.helper.DbModelToSchema(&m))
	}
	return &core.LongTermMemoryArray{
		Memories: memories,
		Total:    int(total),
	}, nil
}

// GetById implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) GetById(ctx context.Context, chatId string, memoryId string) (*core.LongTermMemory, error) {
	dbMemory, err := l.G().
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return l.helper.DbModelToSchema(&dbMemory), nil
}

// GetScored implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) GetScored(
	ctx context.Context,
	chatId string,
	memoriesIds []string,
) ([]*core.ScoredMemory, error) {
	var scoredMemories []*core.ScoredMemory
	if len(memoriesIds) == 0 {
		return scoredMemories, nil
	}
	dbMemories, err := l.G().Where("id IN ?", memoriesIds).Find(ctx)
	if err != nil {
		slog.Error("failed to get memories", "error", err)
		return nil, err
	}
	var rawMemories []core.LongTermMemory
	for _, m := range dbMemories {
		rawMemories = append(rawMemories, *l.helper.DbModelToSchema(&m))
	}

	maxAge, maxAccessCount, err := l.helper.GetMaxCounts(rawMemories)
	if err != nil {
		slog.Error("failed to get max counts", "error", err)
		return nil, err
	}
	slog.Info("max age", "maxAge", maxAge, "maxAccessCount", maxAccessCount)
	now := time.Now().Unix()
	for _, memory := range rawMemories {
		score, err := l.helper.CalculateScore(
			memory,
			maxAge,
			maxAccessCount,
			now,
		)
		if err != nil {
			return nil, err
		}
		scoredMemories = append(
			scoredMemories, &core.ScoredMemory{
				Id:             memory.Id,
				Text:           memory.Memory,
				Score:          score,
				MemoryType:     core.LongTerm,
				CreatedAt:      memory.CreatedAt,
				RelatedContext: memory.RelatedContext,
			},
		)
	}
	return scoredMemories, nil
}

// RegisterUsage implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) RegisterUsage(ctx context.Context, chatId string, memoryId string) error {
	if _, err := l.G().
		Where("chat_id = ? AND id = ?", chatId, memoryId).
		Update(ctx, "access_count", gorm.Expr("access_count + ?", 1)); err != nil {
		return err
	}
	return nil
}

// DeactivateAll implements [repository.LongTermMemoryRepository]
func (l *LongTermMemoryRepository) DeactivateAll(ctx context.Context, chatId string) error {
	if _, err := l.G().
		Where("chat_id = ?", chatId).
		Update(ctx, "active", false); err != nil {
		return err
	}
	return nil
}

var _ repository.LongTermMemoryRepository = (*LongTermMemoryRepository)(nil)
