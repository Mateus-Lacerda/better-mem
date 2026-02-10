package better_mem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mateus-Lacerda/better-mem/pkg/core"
)

// TODO: Add oauth
type BetterMemClient struct {
	baseUrl    string
	httpClient *http.Client
}

func (c *BetterMemClient) CreateChat(externalId string) error {
	req := core.NewChat{ExternalId: externalId}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/chat", c.baseUrl),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"error creating chat: %s - %s",
			resp.Status,
			string(bodyBytes),
		)
	}
	return nil
}

func (c *BetterMemClient) SendMessage(
	chatId string,
	message string,
	relatedContext []core.MessageRelatedContext,
) error {
	req := core.NewMessage{
		ChatId:         chatId,
		Message:        message,
		RelatedContext: relatedContext,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	send := func() (*http.Response, error) {
		return c.httpClient.Post(
			fmt.Sprintf("%s/message", c.baseUrl),
			"application/json",
			bytes.NewBuffer(body),
		)
	}

	resp, err := send()
	if resp.StatusCode == http.StatusBadRequest {
		c.CreateChat(chatId)
		resp, err = send()
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"error sending message: %s - %s",
			resp.Status,
			string(bodyBytes),
		)
	}

	return nil
}

func (c *BetterMemClient) FetchMemories(
	chatID string, req core.MemoryFetchRequest,
) ([]core.ScoredMemory, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/memory/chat/%s/fetch", c.baseUrl, chatID),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"error fetching memories: %s - %s",
			resp.Status,
			string(bodyBytes),
		)
	}

	var memories []core.ScoredMemory
	if err := json.NewDecoder(resp.Body).Decode(&memories); err != nil {
		return nil, err
	}

	return memories, nil
}
