package ollama

import (
	"github.com/Mateus-Lacerda/better-mem/internal/llm"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type OllamaProvider struct {
	BaseUrl string
	Model   string
}

func NewLLMProvider(baseUrl, model string) *OllamaProvider {
	provider := OllamaProvider{
		BaseUrl: baseUrl,
		Model:   model,
	}
	if provider.TestProvider() != nil {
		return nil
	}
	return &provider
}

// GetCompletion implements [llm.LLMProvider].
func (o OllamaProvider) GetCompletion(prompt string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"model":  o.Model,
		"prompt": prompt,
		"stream": false,
	})
	res, err := http.Post(o.BaseUrl+"/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(resBody, &result); err != nil {
		return "", err
	}
	return result.Response, nil
}

// TestProvider implements [llm.LLMProvider].
func (o OllamaProvider) TestProvider() error {
	res, err := http.Get(o.BaseUrl)
	if err != nil {
		slog.Error("TestProvider", "err", err)
		return err
	}
	if res.StatusCode > 299 {
		return llm.LLMProviderError(o)
	}
	return nil
}

var _ llm.LLMProvider = OllamaProvider{}
