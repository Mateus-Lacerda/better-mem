package service

import (
	"better-mem/internal/config"
	"better-mem/internal/core"
	"better-mem/internal/task"
	"log/slog"

	"github.com/hibiken/asynq"
)

func AddMessage(chatId, message string, relatedContext []core.MessageRelatedContext) error {
	// Debug related context
	for _, context := range relatedContext {
		slog.Info(
			"related context",
			slog.String("message", context.Context),
			slog.String("role", context.User),
		)
	}
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Database.RedisAddress})

	defer client.Close()

	task, err := task.NewClassifyMessageTask(chatId, message, relatedContext)
	if err != nil {
		return err
	}
	info, err := client.Enqueue(task)
	if err != nil {
		return err
	}
	slog.Info(
		"message added to queue",
		slog.String("type", task.Type()),
		slog.String("queue", info.Queue),
		slog.Int("retry", int(info.State)),
	)

	return nil
}
