//go:build server

package main

import (
	"better-mem/internal/config"
	"better-mem/internal/task"
	"better-mem/internal/task/handler"
	"log/slog"

	"better-mem/internal/database/mongo"
	"better-mem/internal/database/mongo/repository"
	"better-mem/internal/database/mongo/uow"
	vectorRepo "better-mem/internal/database/qdrant/repository"
	contracts "better-mem/internal/repository"
	vectorContracts "better-mem/internal/repository/vector"
	uowContracts "better-mem/internal/uow"

	"github.com/hibiken/asynq"
)

func startScheduler() {
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: config.Database.RedisAddress},
		nil,
	)

	scheduler.Register(
		config.MemoryManagement.ManageSTMemoryTaskPeriod,
		task.NewManageShortTermMemoryTask(),
	)
	if err := scheduler.Run(); err != nil {
		slog.Error("failed to run scheduler", "err", err)
		return
	}
}

func startConsumer(
	messageHandler *handler.MessageTaskHandler,
	manageShortTermMemoryHandler *handler.MemoryManagementHandler,
) {

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Database.RedisAddress},
		asynq.Config{
			Concurrency: config.Worker.Concurrency,
			Queues:      queues,
		},
	)

	// Mux
	mux := asynq.NewServeMux()
	mux.HandleFunc(
		task.ClassifyMessageTaskName,
		messageHandler.HandleClassifyMemoryTask,
	)
	mux.HandleFunc(
		task.StoreLongTermMemoryTaskName,
		messageHandler.HandleStoreLongTermMemoryTask,
	)
	mux.HandleFunc(
		task.StoreShortTermMemoryTaskName,
		messageHandler.HandleStoreShortTermMemoryTask,
	)
	mux.HandleFunc(
		task.ManageMemoryTaskName,
		manageShortTermMemoryHandler.HandleManageMemory,
	)

	if err := server.Run(mux); err != nil {
		slog.Error("failed to run server", "err", err)
		return
	}
}

func getRepositories() (
	contracts.ChatRepository,
	contracts.LongTermMemoryRepository,
	contracts.ShortTermMemoryRepository,
	vectorContracts.MemoryVectorRepository,
	uowContracts.UnitOfWork[int, any],
) {
	chatRepository := repository.NewChatRepository()
	longTermMemoryRepository := repository.NewLongTermMemoryRepository()
	shortTermMemoryRepository := repository.NewShortTermMemoryRepository()
	memoryVectorRepository := vectorRepo.NewMemoryRepository()
	mongoIntUow := uow.NewUnitOfWork[int](mongo.GetMongoClient())
	return chatRepository, longTermMemoryRepository, shortTermMemoryRepository, memoryVectorRepository, mongoIntUow
}
