package core

import "time"

type NewLongTermMemory struct {
	Memory         string                  `json:"memory"`
	ChatId         string                  `json:"chat_id"`
	AccessCount    int                     `json:"access_count"`
	CreatedAt      time.Time               `json:"created_at"`
	Active         bool                    `json:"active"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

type LongTermMemory struct {
	Id             string                  `json:"id" bson:"_id"`
	Memory         string                  `json:"memory"`
	ChatId         string                  `json:"chat_id"`
	AccessCount    int                     `json:"access_count"`
	CreatedAt      time.Time               `json:"created_at"`
	Active         bool                    `json:"active"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

type LongTermMemoryArray struct {
	Memories []*LongTermMemory `json:"memories"`
	Total    int               `json:"total"`
}

type LongTermMemoryModel struct {
	Id             string                  `json:"id" bson:"_id"`
	Memory         string                  `json:"memory"`
	ChatId         string                  `json:"chatid"`
	AccessCount    int                     `json:"accesscount"`
	CreatedAt      time.Time               `json:"createdat"`
	Active         bool                    `json:"active"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}
