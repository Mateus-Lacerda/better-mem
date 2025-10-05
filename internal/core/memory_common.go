package core

type MemoryTypeEnum int

const (
	NoMemory MemoryTypeEnum = iota
	ShortTerm
	LongTerm
)

type ScoredMemory struct {
	Text       string  `json:"text"`
	Score      float32 `json:"score"`
	MemoryType MemoryTypeEnum
}

type MemoryPayload struct {
	ChatId     string         `json:"chat_id"`
	MemoryType MemoryTypeEnum `json:"memory_type"`
	MemoryId   string         `json:"memory_id"`
	Active     bool           `json:"active"`
}

func (m *MemoryPayload) ToMap() map[string]any {
	return map[string]any{
		"chat_id":     m.ChatId,
		"memory_type": int(m.MemoryType),
		"memory_id":   m.MemoryId,
		"active":      m.Active,
	}
}

type MemoryVectorModel struct {
	Id      string        `json:"id"`
	Vectors []float32     `json:"vectors"`
	Payload MemoryPayload `json:"payload"`
}

type ScoredMemoryVector struct {
	Id      string        `json:"id"`
	Vectors []float32     `json:"vectors"`
	Score   float32       `json:"score"`
	Payload MemoryPayload `json:"payload"`
}

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

type MemoryManagementResult struct {
	ChatId    string
	Success   bool
	Promoted  int
	Discarded int
}
