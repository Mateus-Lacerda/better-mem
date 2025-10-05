package task

import (
	"time"
	"better-mem/internal/config"

	"github.com/hibiken/asynq"
)

// Definitions
const (
	ManageMemoryTaskName = "short-term-memory:manage"
)

func NewManageShortTermMemoryTask() *asynq.Task {
	return asynq.NewTask(
		ManageMemoryTaskName,
		nil,
		asynq.MaxRetry(config.WorkerConfig.MaxRetry),
		asynq.Timeout(time.Duration(config.WorkerConfig.Timeout)*time.Second),
	)
}
