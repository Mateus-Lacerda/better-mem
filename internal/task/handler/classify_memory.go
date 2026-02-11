package handler

import (
	"github.com/Mateus-Lacerda/better-mem/internal/config"
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	protos "github.com/Mateus-Lacerda/better-mem/internal/grpc_client"
	"github.com/Mateus-Lacerda/better-mem/internal/service"
	"github.com/Mateus-Lacerda/better-mem/internal/task"
	"context"
	"encoding/json"
	"log/slog"
)

type MessageTaskHandler struct {
	longTermMemoryService    *service.LongTermMemoryService
	shortTermMemoryService   *service.ShortTermMemoryService
	memoryVectorService      *service.MemoryVectorService
	memoryEnhancementService *service.MemoryEnhancementService
}

func NewMessageTaskHandler(
	longTermMemoryService *service.LongTermMemoryService,
	shortTermMemoryService *service.ShortTermMemoryService,
	memoryVectorService *service.MemoryVectorService,
	memoryEnhancementService *service.MemoryEnhancementService,
) *MessageTaskHandler {
	return &MessageTaskHandler{
		longTermMemoryService:    longTermMemoryService,
		shortTermMemoryService:   shortTermMemoryService,
		memoryVectorService:      memoryVectorService,
		memoryEnhancementService: memoryEnhancementService,
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
// TODO: Fix the memory enhancement and remove the debug slogs
func (h *MessageTaskHandler) handleClassifyMemoryTask(
	ctx context.Context,
	payload task.ClassifyMessagePayload,
	enqueueFunc func(string, []byte) error,
) error {
	slog.Info("handleClassifyMemoryTask", "payload", payload)
	hasEnhancementCapabilites := h.memoryEnhancementService.IsWorking()

	labeledMessage, err := protos.Predict(payload.Message, payload.ChatId, !hasEnhancementCapabilites)
	if err != nil {
		slog.Error("Error predicting message", "error", err)
		return err
	}

	if labeledMessage.Label == core.NoMemory {
		return nil
	}

	if hasEnhancementCapabilites {
		originalMessage := labeledMessage.Message
		enhancedMemory := h.memoryEnhancementService.EnhanceMemory(originalMessage)
		embeddings, err := protos.Embed(enhancedMemory)
		slog.Error("embedding error", "err", err)
		if err == nil {
			slog.Info("updating", "labeledMessage", labeledMessage)
			labeledMessage.Message = enhancedMemory
			labeledMessage.MessageEmbedding = embeddings
		}

		labeledMessage.RelatedContext = append(
			labeledMessage.RelatedContext,
			core.MessageRelatedContext{
				Context: originalMessage,
				User:    "user",
			},
		)
	}
	slog.Info("final", "labeledMessage", labeledMessage)
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
				ctx,
				payload.ChatId,
				similarMemory.Payload.MemoryId,
				payload.Message,
				payload.RelatedContext,
			)
		}
		return nil
	}
	labeledMessage.RelatedContext = payload.RelatedContext

	storeMemoryPayload := task.StoreMemoryPayload{
		LabeledMessage: *labeledMessage,
	}
	slog.Info("store", "storeMemoryPayload", storeMemoryPayload)
	payloadBytes, err := json.Marshal(storeMemoryPayload)
	if err != nil {
		return err
	}

	storeTaskName, err := func() (string, error) {
		switch labeledMessage.Label {
		case core.LongTerm:
			return task.StoreLongTermMemoryTaskName, nil
		case core.ShortTerm:
			return task.StoreShortTermMemoryTaskName, nil
		default:
			return "", core.UnexpectedClassificationError
		}
	}()

	if err != nil {
		return err
	}

	if err = enqueueFunc(storeTaskName, payloadBytes); err != nil {
		return err
	}

	return nil
}

// Saves the message as long term memory
func (h *MessageTaskHandler) handleStoreLongTermMemoryTask(
	ctx context.Context,
	payload task.StoreMemoryPayload,
) error {
	createdMemory, err := h.longTermMemoryService.Create(
		ctx,
		payload.Message,
		payload.ChatId,
		payload.RelatedContext,
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
func (h *MessageTaskHandler) handleStoreShortTermMemoryTask(
	ctx context.Context,
	payload task.StoreMemoryPayload,
) error {
	createdMemory, err := h.shortTermMemoryService.Create(
		ctx,
		payload.Message,
		payload.ChatId,
		payload.RelatedContext,
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
