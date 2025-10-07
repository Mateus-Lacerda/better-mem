package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"better-mem/internal/config"
	"better-mem/internal/core"
	protos "better-mem/internal/grpc_client"
	"better-mem/internal/service"
	"better-mem/internal/task"

	"github.com/hibiken/asynq"
)

type MessageTaskHandler struct {
	longTermMemoryService  *service.LongTermMemoryService
	shortTermMemoryService *service.ShortTermMemoryService
	memoryVectorService    *service.MemoryVectorService
}

func NewMessageTaskHandler(
	longTermMemoryService *service.LongTermMemoryService,
	shortTermMemoryService *service.ShortTermMemoryService,
	memoryVectorService *service.MemoryVectorService,
) *MessageTaskHandler {
	return &MessageTaskHandler{
		longTermMemoryService:  longTermMemoryService,
		shortTermMemoryService: shortTermMemoryService,
		memoryVectorService:    memoryVectorService,
	}
}

func (h *MessageTaskHandler) checkForSimilarMemory(
	ctx context.Context,
	chatId string,
	tokens []float32,
) (*core.ScoredMemoryVector, error) {

	similarMemory, err := h.memoryVectorService.SearchMemoryVector(
		ctx,
		chatId,
		tokens,
		1,
		config.MemoryManagement.MemorySimilarityThreshold,
	)
	if similarMemory == nil {
		return nil, nil
	}
	if len(*similarMemory) > 0 {
		return &(*similarMemory)[0], nil
	}

	return nil, err
}

// ClassifyMemoryTaskHandler handles the heaviest task:
// Classify the message type (long term, short term, none)
func (h *MessageTaskHandler) HandleClassifyMemoryTask(
	ctx context.Context, t *asynq.Task,
) error {
	var payload task.ClassifyMessagePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	labeledMessage, err := protos.Predict(payload.Message, payload.ChatId, true)
	if err != nil {
		slog.Error("Error predicting message", "error", err)
		return err
	}

	if labeledMessage.Label == core.NoMemory {
		return nil
	}

	similarMemory, err := h.checkForSimilarMemory(
		ctx, payload.ChatId, labeledMessage.MessageEmbedding,
	)
	if err != nil {
		return err
	}
	if similarMemory != nil {
		slog.Info("Similar memory found", "message", payload.Message)
		if labeledMessage.Label == core.ShortTerm {
			// Merge the memories
			h.shortTermMemoryService.Merge(
				ctx, payload.ChatId, similarMemory.Payload.MemoryId, payload.Message,
			)
		}
		return nil
	}

	storeMemoryPayload := task.StoreMemoryPayload{
		LabeledMessage: *labeledMessage,
	}
	payloadBytes, err := json.Marshal(storeMemoryPayload)
	if err != nil {
		return err
	}
	storeTask := asynq.NewTask(task.StoreLongTermMemoryTaskName, payloadBytes)

	switch labeledMessage.Label {
	case core.LongTerm:
		if err := h.HandleStoreLongTermMemoryTask(ctx, storeTask); err != nil {
			slog.Error("ClassifyMemoryTaskHandler", "error", err)
			return err
		}
	case core.ShortTerm:
		if err := h.HandleStoreShortTermMemoryTask(ctx, storeTask); err != nil {
			slog.Error("ClassifyMemoryTaskHandler", "error", err)
			return err
		}
	}

	slog.Info("ClassifyMemoryTaskHandler", "label", labeledMessage.Label)
	return nil
}

// Saves the message as long term memory
func (h *MessageTaskHandler) HandleStoreLongTermMemoryTask(
	ctx context.Context, t *asynq.Task,
) error {
	var payload task.StoreMemoryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		slog.Error("HandleStoreLongTermMemoryTask", "error", err)
		return err
	}
	createdMemory, err := h.longTermMemoryService.Create(
		ctx,
		payload.Message,
		payload.ChatId,
	)
	if err != nil {
		slog.Error("HandleStoreLongTermMemoryTask", "error", err)
		return err
	}
	if err := h.memoryVectorService.CreateMemoryVector(
		ctx,
		payload.ChatId,
		payload.MessageEmbedding,
		core.LongTerm,
		createdMemory.Id,
	); err != nil {
		slog.Error("HandleStoreLongTermMemoryTask", "error", err)
		return err
	}
	return nil
}

// Saves the message as short term memory
func (h *MessageTaskHandler) HandleStoreShortTermMemoryTask(
	ctx context.Context, t *asynq.Task,
) error {
	var payload task.StoreMemoryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		slog.Error("HandleStoreShortTermMemoryTask", "error", err)
		return err
	}
	createdMemory, err := h.shortTermMemoryService.Create(
		ctx,
		payload.Message,
		payload.ChatId,
	)
	if err != nil {
		slog.Error("HandleStoreShortTermMemoryTask", "error", err)
		return err
	}
	if err := h.memoryVectorService.CreateMemoryVector(
		ctx,
		payload.ChatId,
		payload.MessageEmbedding,
		core.ShortTerm,
		createdMemory.Id,
	); err != nil {
		slog.Error("HandleStoreShortTermMemoryTask", "error", err)
		return err
	}
	return nil
}
