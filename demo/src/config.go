package demo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Name                  string  `json:"name"`
	ChatID                string  `json:"chat_id"`
	OpenAIKey             string  `json:"openai_key"`
	APIBaseURL            string  `json:"api_base_url"`
	Limit                 int     `json:"limit"`
	VectorSearchLimit     int     `json:"vector_search_limit"`
	VectorSearchThreshold float32 `json:"vector_search_threshold"`
	LongTermThreshold     float32 `json:"long_term_threshold"`
	ChatHistoryBuffer     int     `json:"chat_history_buffer"`
	Model                 string  `json:"model"`
}

func DefaultConfig() *Config {
	return &Config{
		APIBaseURL:            "http://localhost:8080/api/v1",
		Limit:                 2,
		VectorSearchLimit:     10,
		VectorSearchThreshold: 0.6,
		LongTermThreshold:     0.8,
		ChatHistoryBuffer:     20,
		Model:                 "gpt-4.1-mini",
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func CreateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}

func (c *Config) Print() {
	fmt.Printf("\n=== Configurações Atuais ===\n")
	fmt.Printf("Nome: %s\n", c.Name)
	fmt.Printf("Chat ID: %s\n", c.ChatID)
	fmt.Printf("API Base URL: %s\n", c.APIBaseURL)
	fmt.Printf("Limite de Memórias: %d\n", c.Limit)
	fmt.Printf("Limite de Busca Vetorial: %d\n", c.VectorSearchLimit)
	fmt.Printf("Threshold de Busca Vetorial: %.2f\n", c.VectorSearchThreshold)
	fmt.Printf("Threshold de Longo Prazo: %.2f\n", c.LongTermThreshold)
	fmt.Printf("Buffer de Histórico: %d\n", c.ChatHistoryBuffer)
	fmt.Printf("Modelo: %s\n", c.Model)
	fmt.Printf("===========================\n\n")
}
