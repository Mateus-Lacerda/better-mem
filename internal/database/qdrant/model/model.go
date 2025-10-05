package model

import "better-mem/internal/core"

type MemoryPayload struct {
	MemoryType core.MemoryTypeEnum `json:"memory_type"`
	MemoryId   string              `json:"memory_id"`
	Active     bool                `json:"active"`
}
