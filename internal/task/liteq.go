//go:build local

package task

import (
	"better-mem/internal/core"
	"better-mem/internal/database/sqlite"
	"context"

	"github.com/khepin/liteq"
)

func NewClassifyMessageTask(
	chatId, message string, relatedContext []core.MessageRelatedContext,
) error {
	db, err := sqlite.GetDb().DB()
	if err != nil {
		return err
	}
	jqueue := liteq.New(db)
	payload, err := getClassifiyMessageTaskPayload(chatId, message, relatedContext)
	return jqueue.QueueJob(
		context.Background(),
		liteq.QueueJobParams{
			Queue: ClassifyMessageTaskName,
			Job:   string(payload),
		},
	)
}
