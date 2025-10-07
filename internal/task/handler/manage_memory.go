package handler

import (
	"context"
	"log/slog"
	"sync"
	"better-mem/internal/config"
	"better-mem/internal/core"
	"better-mem/internal/service"

	"github.com/hibiken/asynq"
)

type MemoryManagementHandler struct {
	chatService             *service.ChatService
	memoryManagementService *service.MemoryManagementService
	completed               chan core.MemoryManagementResult
	wgCompleted             sync.WaitGroup
}

func NewMemoryManagementHandler(
	chatService *service.ChatService,
	memoryManagementService *service.MemoryManagementService,
) *MemoryManagementHandler {

	handler := &MemoryManagementHandler{
		chatService:             chatService,
		memoryManagementService: memoryManagementService,
		completed: make(
			chan core.MemoryManagementResult,
			config.MemoryManagement.MaxSimultaneousTasks,
		),
		wgCompleted: sync.WaitGroup{},
	}
	handler.wgCompleted.Add(1)
	go handler.HandleCompleted()
	return handler
}

func (m *MemoryManagementHandler) ManageShortTermMemory(
	ctx context.Context, chatId string,
) {
	var success bool
	var promoted, deactivated int
	deactivated, err := m.memoryManagementService.FindAndDeactivate(
		ctx,
		chatId,
		config.MemoryManagement.STValConfig.AgeLimitHours,
		config.MemoryManagement.STValConfig.MinimalRelevancyForDiscard,
	)
	if err != nil {
		slog.Error("error deactivating short term memories", "error", err)
		success = false
	}
	promoted, err = m.memoryManagementService.FindAndPromote(
		ctx,
		chatId,
		config.MemoryManagement.STValConfig.MinimalRelevancyForPromotion,
		config.MemoryManagement.STValConfig.LongTermThreshold,
	)
	m.completed <- core.MemoryManagementResult{
		ChatId:    chatId,
		Success:   success,
		Promoted:  promoted,
		Discarded: deactivated,
	}
}

func (m *MemoryManagementHandler) HandleManageMemory(
	ctx context.Context, t *asynq.Task,
) error {
	chats, err := m.chatService.GetAll(ctx)
	if err != nil {
		slog.Error("error getting chats", "error", err)
		return err
	}
	slog.Info("managing memory", "chats", len(chats))
	if len(chats) == 0 {
		return nil
	}
	jobs := make(chan string, len(chats))

	var wgWorkers sync.WaitGroup

	for range config.MemoryManagement.MaxSimultaneousTasks {
		wgWorkers.Go(func() {
			for chatId := range jobs {
				m.ManageShortTermMemory(ctx, chatId)
			}
		})
	}
	for _, chat := range chats {
		jobs <- chat.ExternalId
	}
	close(jobs)

	wgWorkers.Wait()
	return nil
}

func (m *MemoryManagementHandler) HandleCompleted() {
	defer m.wgCompleted.Done()
	for completed := range m.completed {
		slog.Info(
			"memory management completed",
			"chat_id", completed.ChatId,
			"success", completed.Success,
			"promoted", completed.Promoted,
			"discarded", completed.Discarded,
		)
	}
	slog.Info("memory management completed")
}

func (m *MemoryManagementHandler) Stop() {
	close(m.completed)
	m.wgCompleted.Wait()
}
