package model

import "github.com/Mateus-Lacerda/better-mem/pkg/core"

type MemoryPayload struct {
	MemoryType core.MemoryTypeEnum `json:"memory_type"`
	MemoryId   string              `json:"memory_id"`
	Active     bool                `json:"active"`
}
