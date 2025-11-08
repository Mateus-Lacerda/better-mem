package v1

import (
	"strconv"
	"better-mem/internal/core"
	"better-mem/internal/service"

	"github.com/gin-gonic/gin"
)


type MemoryHandler struct {
	shortTermMemoryService *service.ShortTermMemoryService
	longTermMemoryService *service.LongTermMemoryService
	memoryService *service.MemoryService
}

func NewMemoryHandler(
	shortTermMemoryService *service.ShortTermMemoryService,
	longTermMemoryService *service.LongTermMemoryService,
	memoryService *service.MemoryService,
) *MemoryHandler {
	return &MemoryHandler{
		shortTermMemoryService: shortTermMemoryService,
		longTermMemoryService: longTermMemoryService,
		memoryService: memoryService,
	}
}

// @Summary Fetch Memories
// @Description Fetch memories for a given chat.
// @Tags memories
// @Accept json
// @Produce json
// @Param chat_id path string true "Chat ID"
// @Param request body core.MemoryFetchRequest true "Fetch Memories Request"
// @Success 200 {object} []core.ScoredMemory
// @Router /memory/chat/{chat_id}/fetch [post]
func (h *MemoryHandler) FetchMemories(context *gin.Context) {
	var request core.MemoryFetchRequest
	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chat_id := context.Param("chat_id")
	if chat_id == "" {
		context.JSON(400, gin.H{"error": "Invalid chat id"})
		return
	}
	memories, err := h.memoryService.Fetch(
		context, 
		chat_id, 
		request.Text,
		request.Limit, 
		request.VectorSearchLimit,
		request.VectorSearchThreshold,
		request.LongTermThreshold,
	)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, memories)
}


// @Summary Get Long Term Memories
// @Description Get long term memories for a given chat.
// @Tags memories
// @Accept json
// @Produce json
// @Param chat_id path string true "Chat ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} []core.LongTermMemory
// @Router /memory/long-term/chat/{chat_id} [get]
func (h *MemoryHandler) GetLongTermMemories(context *gin.Context) {
	chat_id := context.Param("chat_id")
	limitStr := context.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		context.JSON(400, gin.H{"error": "Invalid limit"})
		return
	}
	offsetStr := context.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		context.JSON(400, gin.H{"error": "Invalid offset"})
		return
	}
	memories, err := h.longTermMemoryService.GetByChatId(
		context, chat_id, limit, offset,
	)
	if err == core.ChatNotFound {
		context.JSON(404, gin.H{"error": "Chat not found"})
		return
	}
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, memories)
}


// @Summary Get Short Term Memories
// @Description Get short term memories for a given chat.
// @Tags memories
// @Accept json
// @Produce json
// @Param chat_id path string true "Chat ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} []core.ShortTermMemory
// @Router /memory/short-term/chat/{chat_id} [get]
func (h *MemoryHandler) GetShortTermMemories(context *gin.Context) {
	chat_id := context.Param("chat_id")
	limitStr := context.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		context.JSON(400, gin.H{"error": "Invalid limit"})
		return
	}
	offsetStr := context.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		context.JSON(400, gin.H{"error": "Invalid offset"})
		return
	}
	memories, err := h.shortTermMemoryService.GetByChatId(
		context, chat_id, limit, offset,
	)
	if err == core.ChatNotFound {
		context.JSON(404, gin.H{"error": "Chat not found"})
		return
	}
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, memories)
}
