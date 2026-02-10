package main

import (
	"github.com/Mateus-Lacerda/better-mem/internal/config"
	"github.com/Mateus-Lacerda/better-mem/internal/llm/ollama"
	"github.com/Mateus-Lacerda/better-mem/internal/service"
	"github.com/Mateus-Lacerda/better-mem/internal/task/handler"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var queues = map[string]int{
	"critical": 6,
	"default":  3,
	"low":      1,
}

func startServer() {
	// Providers
	llmProvider := ollama.NewLLMProvider(config.Llm.BaseUrl, config.Llm.Model)

	// Repositories
	chatRepository,
		longTermMemoryRepository,
		shortTermMemoryRepository,
		memoryVectorRepository, uow := getRepositories()

	// Services
	longTermMemoryService := service.NewLongTermMemoryService(longTermMemoryRepository, chatRepository)
	shortTermMemoryService := service.NewShortTermMemoryService(shortTermMemoryRepository, chatRepository)
	chatService := service.NewChatService(chatRepository)
	memoryVectorService := service.NewMemoryVectorService(memoryVectorRepository)
	memoryManagementService := service.NewMemoryManagementService(uow)
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
	startConsumer(messageHandler, manageShortTermMemoryHandler)
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
