package task

import (
	"better-mem/internal/config"

	"github.com/hibiken/asynq"
)

func getAsynqClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: config.Database.RedisAddress})
}

func Enqueue(task *asynq.Task) (*asynq.TaskInfo, error) {
	client := getAsynqClient()
	defer client.Close()
	return client.Enqueue(task)
}
