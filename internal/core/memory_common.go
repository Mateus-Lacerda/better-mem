package core

import "time"

type MemoryTypeEnum int

const (
	NoMemory MemoryTypeEnum = iota
	ShortTerm
	LongTerm
)

// A memory with a score
type ScoredMemory struct {
	// The id of the memory
	Id string `json:"id"`
	// The text that was used to generate the memory
	Text string `json:"text"`
	// The score of the memory
	Score float32 `json:"score"`
	// The date the memory was created
	CreatedAt time.Time `json:"created_at"`
	// The type of the memory
	MemoryType MemoryTypeEnum `json:"memory_type"`
	// Context that might be related to the memory
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

// Payload for the memory that is stored in the vector database
type MemoryPayload struct {
	ChatId     string         `json:"chat_id"`
	MemoryType MemoryTypeEnum `json:"memory_type"`
	MemoryId   string         `json:"memory_id"`
	Active     bool           `json:"active"`
}

// ToMap converts the memory payload to a map that can
// be used to store in the vector database
func (m *MemoryPayload) ToMap() map[string]any {
	return map[string]any{
		"chat_id":     m.ChatId,
		"memory_type": int(m.MemoryType),
		"memory_id":   m.MemoryId,
		"active":      m.Active,
	}
}

// Vector model for the memory
type MemoryVectorModel struct {
	Id      string        `json:"id"`
	Vectors []float32     `json:"vectors"`
	Payload MemoryPayload `json:"payload"`
}

// A memory vector with a score
type ScoredMemoryVector struct {
	Id      string        `json:"id"`
	Vectors []float32     `json:"vectors"`
	Score   float32       `json:"score"`
	Payload MemoryPayload `json:"payload"`
}

// Request schema for fetching memories
type MemoryFetchRequest struct {
	// Text to be searched
	Text string `json:"text" example:"I love smart LLMs"`
	// Max number of memories to be returned (Default: 2)
	Limit int `json:"limit" example:"2"`
	// Max number of memories to be returned from vector search (Default: 10)
	VectorSearchLimit int `json:"vector_search_limit" example:"10"`
	// Min score to considerate a memory (Default: 0.6)
	VectorSearchThreshold float32 `json:"vector_search_threshold" example:"0.4"`
	// Specific threshold for long term memories (Default: 0.8)
	LongTermThreshold float32 `json:"long_term_threshold" example:"0.6"`
}

// Result of the memory management
type MemoryManagementResult struct {
	ChatId    string
	Success   bool
	Promoted  int
	Discarded int
}
