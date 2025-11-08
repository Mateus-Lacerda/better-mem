package repository

import (
	"context"
	"better-mem/internal/core"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *core.Chat) error
	GetAll(ctx context.Context) ([]*core.Chat, error)
	GetByExternalID(ctx context.Context, externalID string) (*string, error)
}


type MockChatRepository struct {
	Chats []*core.Chat
}

func NewMockChatRepository() *MockChatRepository {
	mockChats := []*core.Chat{
		{ExternalId: "1"},
		{ExternalId: "2"},
		{ExternalId: "3"},
		{ExternalId: "4"},
	}
	return &MockChatRepository{Chats: mockChats}
}

func (m *MockChatRepository) Create(ctx context.Context, chat *core.Chat) error {
	m.Chats = append(m.Chats, chat)
	return nil
}

func (m *MockChatRepository) GetAll(ctx context.Context) ([]*core.Chat, error) {
	return m.Chats, nil
}
