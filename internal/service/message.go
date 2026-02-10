package service

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"github.com/Mateus-Lacerda/better-mem/internal/task"
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
