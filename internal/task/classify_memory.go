package task

import (
	"encoding/json"
	"time"
	"better-mem/internal/config"
	"better-mem/internal/core"

	"github.com/hibiken/asynq"
)

// Definitions
const (
	// ClassifyMessageTaskName is the name of the classify memory task.
	ClassifyMessageTaskName = "message:classify"
	// StoreLongTermMemoryTaskName is the name of the store long term memory task.
	StoreLongTermMemoryTaskName = "long-term-memory:store"
	// StoreShortTermMemoryTaskName is the name of the store short term memory task.
	StoreShortTermMemoryTaskName = "short-term-memory:store"
)

type ClassifyMessagePayload struct {
	core.NewMessage `json:"embedded"`
}

type StoreMemoryPayload struct {
	core.LabeledMessage `json:"embedded"`
}

func NewClassifyMessageTask(
	chatId, message string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(
		ClassifyMessagePayload{
			core.NewMessage{ChatId: chatId, Message: message},
		},
	)
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

func NewStoreLongTermMemoryTask(
	chatId,
	message string,
	memoryType core.MemoryTypeEnum,
	messageEmbedding []float32,
) (*asynq.Task, error) {
	payload, err := json.Marshal(
		StoreMemoryPayload{
			core.LabeledMessage{
				NewMessage:       core.NewMessage{ChatId: chatId, Message: message},
				Label:            memoryType,
				MessageEmbedding: messageEmbedding,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		StoreLongTermMemoryTaskName,
		payload,
		asynq.MaxRetry(config.Worker.MaxRetry),
		asynq.Timeout(time.Duration(config.Worker.Timeout)*time.Second),
	), nil
}

func NewStoreShortTermMemoryTask(
	chatId,
	message string,
	memoryType core.MemoryTypeEnum,
	messageEmbedding []float32,
) (*asynq.Task, error) {
	payload, err := json.Marshal(
		StoreMemoryPayload{
			core.LabeledMessage{
				NewMessage:       core.NewMessage{ChatId: chatId, Message: message},
				Label:            memoryType,
				MessageEmbedding: messageEmbedding,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		StoreShortTermMemoryTaskName,
		payload,
		asynq.MaxRetry(config.Worker.MaxRetry),
		asynq.Timeout(time.Duration(config.Worker.Timeout)*time.Second),
	), nil
}
