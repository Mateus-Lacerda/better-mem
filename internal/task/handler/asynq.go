//go:build server

package handler

import (
	"better-mem/internal/config"
	"better-mem/internal/task"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

func enqueueFunc(taskName string, payloadBytes []byte) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Database.RedisAddress})
	defer client.Close()
	if _, err := client.Enqueue(asynq.NewTask(taskName, payloadBytes)); err != nil {
		return err
	}
	return nil
}

func (h *MessageTaskHandler) HandleClassifyMemoryTask(ctx context.Context, t *asynq.Task) error {
	var payload task.ClassifyMessagePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return h.handleClassifyMemoryTask(ctx, payload, enqueueFunc)
}

func (h *MessageTaskHandler) HandleStoreLongTermMemoryTask(
	ctx context.Context, t *asynq.Task,
) error {
	var payload task.StoreMemoryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return h.handleStoreLongTermMemoryTask(ctx, payload)
}

func (h *MessageTaskHandler) HandleStoreShortTermMemoryTask(
	ctx context.Context, t *asynq.Task,
) error {
	var payload task.StoreMemoryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return h.handleStoreLongTermMemoryTask(ctx, payload)
}

func (m *MemoryManagementHandler) HandleManageMemory(
	ctx context.Context, t *asynq.Task,
) error {
	return m.handleManageMemory(ctx)
}
