package service

import (
	"better-mem/internal/core"
	"better-mem/internal/task"
	"log/slog"
)

func AddMessage(chatId, message string, relatedContext []core.MessageRelatedContext) error {
	err := task.NewClassifyMessageTask(chatId, message, relatedContext)
	if err != nil {
		return err
	}
	slog.Info("message added to queue")

	return nil
}
