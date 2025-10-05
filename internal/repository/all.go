package repository

type AllRepositories struct {
	Chat            ChatRepository
	ShortTermMemory ShortTermMemoryRepository
	LongTermMemory  LongTermMemoryRepository
}
