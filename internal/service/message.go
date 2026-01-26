package service

import (
	"better-mem/internal/core"
	"better-mem/internal/task"
	"log/slog"
)

func AddMessage(chatId, message string, relatedContext []core.MessageRelatedContext) error {
	newTask, err := task.NewClassifyMessageTask(chatId, message, relatedContext)
	if err != nil {
		return err
	}
	info, err := task.Enqueue(newTask)
	if err != nil {
		return err
	}
	slog.Info(
		"message added to queue",
		slog.String("type", newTask.Type()),
		slog.String("queue", info.Queue),
		slog.Int("retry", int(info.State)),
	)

	return nil
}
