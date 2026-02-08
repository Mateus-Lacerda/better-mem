//go:build local

package handler

import (
	"better-mem/internal/database/sqlite"
	"better-mem/internal/task"
	"context"
	"encoding/json"

	"github.com/khepin/liteq"
)

func enqueueFunc(taskName string, payloadBytes []byte) error {
	// TODO: Do not create a new instance every time
	db, err := sqlite.GetDb().DB()
	if err != nil {
		return err
	}
	jqueue := liteq.New(db)
	return jqueue.QueueJob(
		context.Background(),
		liteq.QueueJobParams{
			Queue: taskName,
			Job:   string(payloadBytes),
		},
	)
}

func (h *MessageTaskHandler) HandleClassifyMemoryTask(
	ctx context.Context, job *liteq.Job,
) error {
	var payload task.ClassifyMessagePayload
	if err := json.Unmarshal([]byte(job.Job), &payload); err != nil {
		return err
	}
	return h.handleClassifyMemoryTask(ctx, payload, enqueueFunc)
}

func (h *MessageTaskHandler) HandleStoreLongTermMemoryTask(
	ctx context.Context, job *liteq.Job,
) error {
	var payload task.StoreMemoryPayload
	if err := json.Unmarshal([]byte(job.Job), &payload); err != nil {
		return err
	}
	return h.handleStoreLongTermMemoryTask(ctx, payload)
}

func (h *MessageTaskHandler) HandleStoreShortTermMemoryTask(
	ctx context.Context, job *liteq.Job,
) error {
	var payload task.StoreMemoryPayload
	if err := json.Unmarshal([]byte(job.Job), &payload); err != nil {
		return err
	}
	return h.handleStoreLongTermMemoryTask(ctx, payload)
}

func (m *MemoryManagementHandler) HandleManageMemory(
	ctx context.Context, job *liteq.Job,
) error {
	return m.handleManageMemory(ctx)
}
