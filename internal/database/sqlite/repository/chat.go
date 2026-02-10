package repository

import (
	"github.com/Mateus-Lacerda/better-mem/pkg/core"
	"context"
	"errors"
	"log/slog"

	"github.com/Mateus-Lacerda/better-mem/internal/database/sqlite"
	"github.com/Mateus-Lacerda/better-mem/internal/repository"

	"gorm.io/gorm"
)

type ChatRepository struct {
	*gorm.DB
}

func NewChatRepository() *ChatRepository {
	db := sqlite.GetDb()
	return &ChatRepository{
		DB: db,
	}
}

func ChatRepositoryWithTransaction(db *gorm.DB) *ChatRepository {
	return &ChatRepository{
		DB: db,
	}
}

func (r *ChatRepository) G() gorm.Interface[sqlite.Chat] {
	return gorm.G[sqlite.Chat](r.DB)
}

// Create implements repository.ChatRepository.
func (r *ChatRepository) Create(ctx context.Context, chat *core.NewChat) error {
	dbChat := sqlite.Chat{ExternalID: chat.ExternalId}
	err := r.G().Create(ctx, &dbChat)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return core.ChatExternalIdAlreadyExists
	}
	if err != nil {
		return err
	}
	slog.Info("Chat created", "id", dbChat.ID)
	return nil
}

// GetByExternalID implements repository.ChatRepository.
func (r *ChatRepository) GetByExternalID(ctx context.Context, externalID string) (*string, error) {
	dbChat, err := r.G().Where("external_id = ?", externalID).First(ctx)
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, core.ChatNotFound
	}
	return &dbChat.ID, nil
}

// GetAll implements repository.ChatRepository.
func (r *ChatRepository) GetAll(ctx context.Context) ([]*core.Chat, error) {
	dbChats, err := r.G().Find(ctx)
	if err != nil {
		slog.Error("Error getting all chats", "error", err)
		return nil, err
	}
	var chats []*core.Chat
	for _, chat := range dbChats {
		chats = append(chats, &core.Chat{
			ExternalId: chat.ExternalID,
			ID:         chat.ID,
		})
	}
	return chats, nil
}

var _ repository.ChatRepository = (*ChatRepository)(nil)
