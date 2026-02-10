//go:build server

package v1

import (
	"github.com/Mateus-Lacerda/better-mem/internal/database/mongo/repository"
	vectorRepo "github.com/Mateus-Lacerda/better-mem/internal/database/qdrant/repository"
	contracts "github.com/Mateus-Lacerda/better-mem/internal/repository"
	vectorContracts "github.com/Mateus-Lacerda/better-mem/internal/repository/vector"
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
