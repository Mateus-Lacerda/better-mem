package service

import (
	"better-mem/internal/llm"
	"fmt"
	"log/slog"
)

type MemoryEnhancementService struct {
	llmProvider llm.LLMProvider
}

func NewMemoryEnhancementService(llmProvider llm.LLMProvider) *MemoryEnhancementService {
	return &MemoryEnhancementService{llmProvider: llmProvider}
}

func (m MemoryEnhancementService) EnhanceMemory(memory string) string {
	prompt := fmt.Sprintf(llm.MemoryEnhancementPrompt, memory)
	enhancedMemory, err := m.llmProvider.GetCompletion(prompt)
	if err != nil {
		slog.Error("Error enhancing memory", "err", err)
		return memory
	}
	slog.Info("Memory enhancement completed", "enhancedMemory", enhancedMemory)
	return enhancedMemory
}
