package repository

import (
	"context"
	"better-mem/internal/core"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *core.NewChat) error
	GetAll(ctx context.Context) ([]*core.Chat, error)
	GetByExternalID(ctx context.Context, externalID string) (*string, error)
}
