package main

import (
	"better-mem/internal/config"
	"better-mem/internal/database/mongo"
	"better-mem/internal/database/mongo/repository"
	"better-mem/internal/database/mongo/uow"
	vectorRepo "better-mem/internal/database/qdrant/repository"
	"better-mem/internal/llm/ollama"
	"better-mem/internal/service"
	"better-mem/internal/task"
	"better-mem/internal/task/handler"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
)

var queues = map[string]int{
	"critical": 6,
	"default":  3,
	"low":      1,
}

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

func startServer() {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Database.RedisAddress},
		asynq.Config{
			Concurrency: config.Worker.Concurrency,
			Queues:      queues,
		},
	)

	// Providers
	llmProvider := ollama.NewLLMProvider(config.Llm.BaseUrl, config.Llm.Model)

	// Repositories
	longTermMemoryRepository := repository.NewLongTermMemoryRepository()
	shortTermMemoryRepository := repository.NewShortTermMemoryRepository()
	chatRepository := repository.NewChatRepository()
	memoryVectorRepository := vectorRepo.NewMemoryRepository()
	mongoIntUow := uow.NewUnitOfWork[int](mongo.GetMongoClient())

	// Services
	longTermMemoryService := service.NewLongTermMemoryService(longTermMemoryRepository, chatRepository)
	shortTermMemoryService := service.NewShortTermMemoryService(shortTermMemoryRepository, chatRepository)
	chatService := service.NewChatService(chatRepository)
	memoryVectorService := service.NewMemoryVectorService(memoryVectorRepository)
	memoryManagementService := service.NewMemoryManagementService(mongoIntUow)
	memoryEnhancementService := service.NewMemoryEnhancementService(llmProvider)

	// Handlers
	messageHandler := handler.NewMessageTaskHandler(
		longTermMemoryService,
		shortTermMemoryService,
		memoryVectorService,
		memoryEnhancementService,
	)
	manageShortTermMemoryHandler := handler.NewMemoryManagementHandler(
		chatService,
		memoryManagementService,
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

func waitForever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		{
			slog.Info("\033[33mInterrupt signal received, exiting in...\033[0m")
			for i := 5; i > 0; i-- {
				slog.Info(fmt.Sprintf("\033[33m%d\033[0m", i))
				time.Sleep(1 * time.Second)
			}
			os.Exit(0)
		}
	}()
	for {
		runtime.Gosched()
	}

}

func main() {
	go startServer()
	go startScheduler()
	waitForever()
}
