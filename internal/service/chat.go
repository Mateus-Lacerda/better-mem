package service

import (
	"context"
	"log/slog"
	"better-mem/internal/core"
	"better-mem/internal/repository"
)

type ChatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) Create(
	ctx context.Context,
	externalId string,
) error {
	chat := &core.Chat{
		ExternalId: externalId,
	}
	err := s.repo.Create(ctx, chat)
	slog.Info("chat created", "chat", chat)
	if err != nil {
		slog.Error("error creating chat", "error", err)
	}
	return err
}

func (s *ChatService) GetByExternalId(ctx context.Context, externalId string) (*string, error) {
	return s.repo.GetByExternalID(ctx, externalId)
}

func (s *ChatService) GetAll(ctx context.Context) ([]*core.Chat, error) {
	return s.repo.GetAll(ctx)
}
