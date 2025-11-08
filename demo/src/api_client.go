package demo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL string
	client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type CreateChatRequest struct {
	ExternalID string `json:"external_id"`
}

type NewMessageRequest struct {
	ChatID         string                  `json:"chat_id"`
	Message        string                  `json:"message"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

type MemoryFetchRequest struct {
	Text                  string  `json:"text"`
	Limit                 int     `json:"limit"`
	VectorSearchLimit     int     `json:"vector_search_limit"`
	VectorSearchThreshold float32 `json:"vector_search_threshold"`
	LongTermThreshold     float32 `json:"long_term_threshold"`
}

type MessageRelatedContext struct {
	Context string `json:"context"`
	User    string `json:"user"`
}

type ScoredMemory struct {
	Text           string                  `json:"text"`
	Score          float32                 `json:"score"`
	MemoryType     int                     `json:"memory_type"`
	CreatedAt      time.Time               `json:"created_at"`
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

func (c *APIClient) CreateChat(externalID string) error {
	req := CreateChatRequest{ExternalID: externalID}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(
		fmt.Sprintf("%s/chat", c.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("falha ao criar chat: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *APIClient) SendMessage(chatID, message string, relatedContext []MessageRelatedContext) error {
	req := NewMessageRequest{
		ChatID:         chatID,
		Message:        message,
		RelatedContext: relatedContext,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(
		fmt.Sprintf("%s/message", c.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("falha ao enviar mensagem: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *APIClient) FetchMemories(chatID string, req MemoryFetchRequest) ([]ScoredMemory, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(
		fmt.Sprintf("%s/memory/chat/%s/fetch", c.baseURL, chatID),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("falha ao buscar memórias: %s - %s", resp.Status, string(bodyBytes))
	}

	var memories []ScoredMemory
	if err := json.NewDecoder(resp.Body).Decode(&memories); err != nil {
		return nil, err
	}

	return memories, nil
}
