package v1

import (
	"better-mem/internal/core"
	"better-mem/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct{
	chatService *service.ChatService
}

func NewMessageHandler(chatService *service.ChatService) *MessageHandler {
	return &MessageHandler{chatService: chatService}
}

type MessageResponse struct {
	Message string `json:"message"`
}

// @Summary Add message
// @Description Sends a message to classification queue
// @Tags message
// @Accept json
// @Produce json
// @Param message body core.NewMessage true "Message"
// @Success 202 {object} MessageResponse
// @Failure 400 {object} any
// @Failure 500 {object} any
// @Router /message [post]
func (h *MessageHandler) AddMessage(context *gin.Context) {
	var m core.NewMessage
	if err := context.BindJSON(&m); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chatId, err := h.chatService.GetByExternalId(context, m.ChatId)
	if err == core.ChatNotFound {
		slog.Info("Chat not found", "chat_id", m.ChatId)
		context.JSON(400, gin.H{"error": "Chat not found"})
		return
	}
	if chatId == nil || err != nil {
		slog.Error("Error getting chat", "error", err)
		context.JSON(500, gin.H{"error": "Error getting chat"})
		return
	}
	if err := service.AddMessage(m.ChatId, m.Message, m.RelatedContext); err != nil {
		slog.Error("Error adding message", "error", err)
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	messageResponse := MessageResponse{Message: "Message accepted"}
	context.Header("Content-Type", "application/json")
	context.JSON(202, messageResponse)
}
