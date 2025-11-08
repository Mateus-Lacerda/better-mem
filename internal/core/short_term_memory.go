package core

import "time"

type NewShortTermMemory struct {
	Memory      string    `json:"memory"`
	ChatId      string    `json:"chat_id"`
	AccessCount int       `json:"access_count"`
	MergeCount  int       `json:"merge_count"`
	Merged      bool      `json:"merged"`
	CreatedAt   time.Time `json:"created_at"`
	Active      bool      `json:"active"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

type ShortTermMemory struct {
	Id             string                  `json:"id" bson:"_id"`
	Memory         string                  `json:"memory"`
	ChatId         string                  `json:"chat_id"`
	AccessCount    int                     `json:"access_count"`
	MergeCount     int                     `json:"merge_count"`
	Merged         bool                    `json:"merged"`
	CreatedAt      time.Time               `json:"created_at"`
	Active         bool                    `json:"active"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

type ShortTermMemoryArray struct {
	Memories []*ShortTermMemory `json:"memories"`
	Total    int                `json:"total"`
}

type ShortTermMemoryModel struct {
	Memory      string    `json:"memory"`
	ChatId      string    `json:"chatid"`
	AccessCount int       `json:"accesscount"`
	MergeCount  int       `json:"mergecount"`
	Merged      bool      `json:"merged"`
	CreatedAt   time.Time `json:"createdat"`
	Active      bool      `json:"active"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}
