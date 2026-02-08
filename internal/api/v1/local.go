//go:build local

package v1

import (
	"better-mem/internal/database/sqlite/repository"
	vectorRepo "better-mem/internal/database/sqlite/repository/vector"
	contracts "better-mem/internal/repository"
	vectorContracts "better-mem/internal/repository/vector"
)

func getRepositories() (
	contracts.ChatRepository,
	contracts.LongTermMemoryRepository,
	contracts.ShortTermMemoryRepository,
	vectorContracts.MemoryVectorRepository,
) {
	chatRepository := repository.NewChatRepository()
	longTermMemoryRepository := repository.NewLongTermMemoryRepository()
	shortTermMemoryRepository := repository.NewShortTermMemoryRepository()
	memoryVectorRepository := vectorRepo.NewMemoryRepository()
	return chatRepository, longTermMemoryRepository, shortTermMemoryRepository, memoryVectorRepository
}
