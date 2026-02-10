package repository

import (
	"context"
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *core.NewChat) error
	GetAll(ctx context.Context) ([]*core.Chat, error)
	GetByExternalID(ctx context.Context, externalID string) (*string, error)
}
