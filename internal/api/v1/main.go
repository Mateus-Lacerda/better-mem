package v1

import (
	"better-mem/internal/database/mongo/repository"
	vectorRepo "better-mem/internal/database/qdrant/repository"
	"better-mem/internal/service"

	"github.com/gin-gonic/gin"
)

// @Summary Health Check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK", "version": "1.0"})
}


func Register(router *gin.Engine) {
	v1Router := router.Group("/api/v1")
	{
		chatRepository := repository.NewChatRepository()
		longTermMemoryRepository := repository.NewLongTermMemoryRepository()
		shortTermMemoryRepository := repository.NewShortTermMemoryRepository()
		memoryVectorRepository := vectorRepo.NewMemoryRepository()
		chatService := service.NewChatService(chatRepository)
		longTermMemoryService := service.NewLongTermMemoryService(longTermMemoryRepository, chatRepository)
		shortTermMemoryService := service.NewShortTermMemoryService(shortTermMemoryRepository, chatRepository)
		memoryService := service.NewMemoryService(
			shortTermMemoryRepository,
			longTermMemoryRepository,
			memoryVectorRepository,
		)
		memoryHandler := NewMemoryHandler(
			shortTermMemoryService,
			longTermMemoryService,
			memoryService,
		)
		chatHandler := NewChatHandler(chatService)
		messageHandler := NewMessageHandler(chatService)

		// Health check
		v1Router.GET("/health", HealthCheck)

		// Memory
		v1Router.GET("/memory/short-term/chat/:chat_id", memoryHandler.GetShortTermMemories)
		v1Router.GET("/memory/long-term/chat/:chat_id", memoryHandler.GetLongTermMemories)
		v1Router.POST("/memory/chat/:chat_id/fetch", memoryHandler.FetchMemories)
	
		// Chat
		v1Router.GET("/chat", chatHandler.GetChats)
		v1Router.POST("/chat", chatHandler.CreateChat)
		
		// Message
		v1Router.POST("/message", messageHandler.AddMessage)
	}
}
