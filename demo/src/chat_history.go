package demo

import (
	"encoding/json"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory struct {
	Messages map[string][]Message `json:"messages"`
}

func LoadChatHistory(path string) (*ChatHistory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ChatHistory{Messages: make(map[string][]Message)}, nil
		}
		return nil, err
	}

	var history ChatHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	if history.Messages == nil {
		history.Messages = make(map[string][]Message)
	}

	return &history, nil
}

func (h *ChatHistory) Save(path string) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (h *ChatHistory) AddMessage(chatID, role, content string) {
	if h.Messages[chatID] == nil {
		h.Messages[chatID] = []Message{}
	}
	h.Messages[chatID] = append(h.Messages[chatID], Message{
		Role:    role,
		Content: content,
	})
}

func (h *ChatHistory) GetMessages(chatID string, buffer int) []Message {
	messages := h.Messages[chatID]
	if len(messages) <= buffer {
		return messages
	}
	return messages[len(messages)-buffer:]
}

