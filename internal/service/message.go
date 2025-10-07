package service

import (
	"log/slog"
	"better-mem/internal/config"
	"better-mem/internal/task"

	"github.com/hibiken/asynq"
)

func AddMessage(chatId, message string) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Database.RedisAddress})

	defer client.Close()

	task, err := task.NewClassifyMessageTask(chatId, message)
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
