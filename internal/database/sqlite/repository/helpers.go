package repository

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	sqlite "github.com/Mateus-Lacerda/better-mem/internal/database/sqlite"
	"time"
)

type ShortTermMemoryHelper struct{}

var ShortTermHelper ShortTermMemoryHelper = ShortTermMemoryHelper{}

func (h *ShortTermMemoryHelper) SchemaToDbModel(
	m *core.NewShortTermMemory,
) *sqlite.ShortTermMemory {
	return &sqlite.ShortTermMemory{
		Memory:      m.Memory,
		ChatID:      m.ChatId,
		AccessCount: m.AccessCount,
		MergeCount:  m.MergeCount,
		Merged:      m.Merged,
		CreatedAt:   m.CreatedAt,
		Active:      m.Active,
	}
}

func (h *ShortTermMemoryHelper) GetMaxCounts(memories []core.ShortTermMemory) (int, int) {
	maxAccessCount := 0
	maxMergeCount := 0
	for _, memory := range memories {
		if memory.AccessCount > maxAccessCount {
			maxAccessCount = memory.AccessCount
		}
		if memory.MergeCount > maxMergeCount {
			maxMergeCount = memory.MergeCount
		}
	}
	return maxAccessCount, maxMergeCount
}

func (h *ShortTermMemoryHelper) GetTemporalScore(now int64, createdAt time.Time) float64 {
	return max(1, time.Since(createdAt).Seconds()/(60*60))
}

func (h *ShortTermMemoryHelper) CalculateScore(
	memory core.ShortTermMemory,
	maxAccessCount int,
	maxMergeCount int,
	now int64,
) (float32, error) {
	relevancyScore :=
		(memory.AccessCount + memory.MergeCount) / max((maxAccessCount)+(maxMergeCount), 1)

	temporalScore := h.GetTemporalScore(now, memory.CreatedAt)
	score := (float32(relevancyScore) + 1/float32(temporalScore)) / 2
	return score, nil
}

func (h *ShortTermMemoryHelper) DbModelToSchema(
	m *sqlite.ShortTermMemory,
) *core.ShortTermMemory {
	var relatedContext []core.MessageRelatedContext
	for _, c := range m.RelatedContext {
		relatedContext = append(
			relatedContext,
			core.MessageRelatedContext{Context: c.Context, User: c.User},
		)
	}
	memory := &core.ShortTermMemory{
		Id:             m.ID,
		Memory:         m.Memory,
		ChatId:         m.ChatID,
		AccessCount:    m.AccessCount,
		MergeCount:     m.MergeCount,
		Merged:         m.Merged,
		CreatedAt:      m.CreatedAt,
		Active:         m.Active,
		RelatedContext: relatedContext,
	}
	return memory
}

type LongTermMemoryHelper struct{}

var LongTermHelper LongTermMemoryHelper = LongTermMemoryHelper{}

func (h *LongTermMemoryHelper) SchemaToDbModel(
	m *core.NewLongTermMemory,
) *sqlite.LongTermMemory {
	return &sqlite.LongTermMemory{
		Memory:      m.Memory,
		ChatID:      m.ChatId,
		AccessCount: m.AccessCount,
		CreatedAt:   m.CreatedAt,
		Active:      m.Active,
	}
}

func (h *LongTermMemoryHelper) GetMaxCounts(
	memories []core.LongTermMemory,
) (int, int, error) {
	maxAge := 0
	maxAccessCount := 0
	for _, memory := range memories {
		age := time.Since(memory.CreatedAt).Seconds() / (60 * 60)
		if int(age) > maxAge {
			maxAge = int(age)
		}
		if memory.AccessCount > maxAccessCount {
			maxAccessCount = memory.AccessCount
		}

	}
	return maxAge, maxAccessCount, nil
}

func (h *LongTermMemoryHelper) CalculateScore(
	memory core.LongTermMemory,
	maxAge int,
	maxAccessCount int,
	now int64,
) (float32, error) {
	age := time.Since(memory.CreatedAt).Seconds() / (60 * 60)
	relevancyScore := float32(memory.AccessCount) / float32(max(maxAccessCount, 1))
	temporalScore := float32(age) / float32(max(maxAge, 1))
	score := (relevancyScore + temporalScore) / 2
	return score, nil
}

func (h *LongTermMemoryHelper) DbModelToSchema(
	m *sqlite.LongTermMemory,
) *core.LongTermMemory {
	var relatedContext []core.MessageRelatedContext
	for _, c := range m.RelatedContext {
		relatedContext = append(
			relatedContext,
			core.MessageRelatedContext{Context: c.Context, User: c.User},
		)
	}
	memory := &core.LongTermMemory{
		Id:             m.ID,
		Memory:         m.Memory,
		ChatId:         m.ChatID,
		AccessCount:    m.AccessCount,
		CreatedAt:      m.CreatedAt,
		Active:         m.Active,
		RelatedContext: relatedContext,
	}
	return memory
}
