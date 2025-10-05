package v1

import (
	"log/slog"
	"better-mem/internal/core"
	"better-mem/internal/service"

	"github.com/gin-gonic/gin"
)


type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// @Summary Create a new chat
// @Description Create a new chat
// @Tags chat
// @Accept json
// @Produce json
// @Param chat body core.Chat true "chat"
// @Success 201 {object} core.Chat
// @Failure 422 {object} string
// @Failure 500 {object} string
// @Router /chat [post]
func (c *ChatHandler) CreateChat(context *gin.Context) {
	var chat core.Chat
	if err := context.BindJSON(&chat); err != nil {
		slog.Error("error binding json", "error", err)
		context.JSON(422, gin.H{"error": err.Error()})
		return
	}
	if err := c.chatService.Create(
		context, chat.ExternalId,
	); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(201, chat)
}

// @Summary Get all chats
// @Description Get all chats
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {array} core.Chat
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /chat [get]
func (c *ChatHandler) GetChats(context *gin.Context) {
	chats, err := c.chatService.GetAll(context)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if chats == nil {
		context.JSON(404, []*core.Chat{})
		return
	}
	context.JSON(200, chats)
}
