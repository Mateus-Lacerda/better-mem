package task

import (
	"better-mem/internal/core"
	"encoding/json"
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

func getClassifiyMessageTaskPayload(
	chatId, message string, relatedContext []core.MessageRelatedContext,
) ([]byte, error) {
	return json.Marshal(
		ClassifyMessagePayload{
			core.NewMessage{ChatId: chatId, Message: message, RelatedContext: relatedContext},
		},
	)
}
