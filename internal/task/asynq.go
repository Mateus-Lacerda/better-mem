//go:build server

package task

import (
	"github.com/Mateus-Lacerda/better-mem/internal/config"
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"time"

	"github.com/hibiken/asynq"
)

func NewClassifyMessageTask(
	chatId, message string, relatedContext []core.MessageRelatedContext,
) (*asynq.Task, error) {
	payload, err := getClassifiyMessageTaskPayload(chatId, message, relatedContext)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		ClassifyMessageTaskName,
		payload,
		asynq.MaxRetry(config.Worker.MaxRetry),
		asynq.Timeout(time.Duration(config.Worker.Timeout)*time.Second),
	), nil
}

func NewManageShortTermMemoryTask() *asynq.Task {
	return asynq.NewTask(
		ManageMemoryTaskName,
		nil,
		asynq.MaxRetry(config.Worker.MaxRetry),
		asynq.Timeout(time.Duration(config.Worker.Timeout)*time.Second),
	)
}
