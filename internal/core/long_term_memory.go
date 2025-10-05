package core

import "time"

type NewLongTermMemory struct {
	Memory      string `json:"memory"`
	ChatId      string `json:"chat_id"`
	AccessCount int    `json:"access_count"`
	CreatedAt   time.Time `json:"created_at"`
	Active      bool   `json:"active"`
}

type LongTermMemory struct {
	Id string `json:"id"`
	NewLongTermMemory
}

type LongTermMemoryArray struct {
	Memories []*LongTermMemory `json:"memories"`
	Total    int              `json:"total"`
}

type LongTermMemoryModel struct {
	Memory string `json:"memory"`
	ChatId string `json:"chatid"`
	AccessCount int    `json:"accesscount"`
	CreatedAt   time.Time `json:"createdat"`
	Active      bool   `json:"active"`
}
