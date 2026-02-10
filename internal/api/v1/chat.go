package v1

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"github.com/Mateus-Lacerda/better-mem/internal/service"
	"errors"
	"log/slog"

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
// @Param chat body core.NewChat true "chat"
// @Success 201 {object} core.Chat
// @Failure 422 {object} string
// @Failure 500 {object} string
// @Router /chat [post]
func (c *ChatHandler) CreateChat(context *gin.Context) {
	var chat core.NewChat
	if err := context.BindJSON(&chat); err != nil {
		slog.Error("error binding json", "error", err)
		context.AbortWithStatusJSON(422, gin.H{"error": err.Error()})
		return
	}
	if err := c.chatService.Create(
		context, chat.ExternalId,
	); err != nil {
		if errors.Is(err, core.ChatExternalIdAlreadyExists) {
			context.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		context.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
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
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /chat [get]
func (c *ChatHandler) GetChats(context *gin.Context) {
	chats, err := c.chatService.GetAll(context)
	if err != nil {
		context.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	if chats == nil {
		context.JSON(200, []*core.Chat{})
		return
	}
	context.JSON(200, chats)
}
