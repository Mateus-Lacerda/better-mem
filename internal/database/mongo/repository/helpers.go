package repository

import (
	"better-mem/internal/core"
	"better-mem/internal/database/mongo"
	"errors"
	"time"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type ShortTermMemoryHelper struct{}

var ShortTermHelper ShortTermMemoryHelper = ShortTermMemoryHelper{}

func (h *ShortTermMemoryHelper) SchemaToDbModel(
	m *core.NewShortTermMemory,
) *mongo.ShortTermMemory {
	return &mongo.ShortTermMemory{
		Memory:      m.Memory,
		ChatID:      m.ChatId,
		AccessCount: m.AccessCount,
		MergeCount:  m.MergeCount,
		Merged:      m.Merged,
		CreatedAt:   m.CreatedAt.String(),
		Active:      m.Active,
	}
}

func (h *ShortTermMemoryHelper) GetMaxCounts(memories []core.ShortTermMemoryModel) (int, int) {
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
	memory core.ShortTermMemoryModel,
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

type LongTermMemoryHelper struct{}

var LongTermHelper LongTermMemoryHelper = LongTermMemoryHelper{}

func (h *LongTermMemoryHelper) SchemaToDbModel(
	m *core.NewLongTermMemory,
) *mongo.LongTermMemory {
	return &mongo.LongTermMemory{
		Memory:      m.Memory,
		ChatID:      m.ChatId,
		AccessCount: m.AccessCount,
		CreatedAt:   m.CreatedAt.String(),
		Active:      m.Active,
	}
}

func (h *LongTermMemoryHelper) GetMaxCounts(
	memories []core.LongTermMemoryModel,
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
	memory core.LongTermMemoryModel,
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

func IsMongoDuplicateKeyError(err error) bool {
	var writeEx mongoDriver.WriteException
	if errors.As(err, &writeEx) {
		for _, e := range writeEx.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}
