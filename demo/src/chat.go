package demo

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	openai "github.com/sashabaranov/go-openai"
)

type ChatSession struct {
	config      *Config
	apiClient   *APIClient
	openaiClient *openai.Client
	history     *ChatHistory
	historyPath string
}

func NewChatSession(config *Config, apiClient *APIClient, history *ChatHistory, historyPath string) *ChatSession {
	return &ChatSession{
		config:       config,
		apiClient:    apiClient,
		openaiClient: openai.NewClient(config.OpenAIKey),
		history:      history,
		historyPath:  historyPath,
	}
}

func (s *ChatSession) Start() error {
	yellow := color.New(color.FgYellow).SprintFunc()
	orange := color.New(color.FgHiRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("\n%s\n", cyan("=== Iniciando Chat ==="))
	fmt.Printf("%s\n", green("Digite 'exit' para sair do chat\n"))

	for {
		prompt := promptui.Prompt{
			Label: "Você",
		}

		userInput, err := prompt.Run()
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(userInput)) == "exit" {
			fmt.Println(cyan("\nEncerrando chat..."))
			break
		}

		if strings.TrimSpace(userInput) == "" {
			continue
		}

		time.Sleep(2 * time.Second)

		memories, err := s.apiClient.FetchMemories(s.config.ChatID, MemoryFetchRequest{
			Text:                  userInput,
			Limit:                 s.config.Limit,
			VectorSearchLimit:     s.config.VectorSearchLimit,
			VectorSearchThreshold: s.config.VectorSearchThreshold,
			LongTermThreshold:     s.config.LongTermThreshold,
		})
		if err != nil {
			fmt.Printf("Error fetching memories from better-mem: %v\n", err)
		}


		// Send last 2 messages as related context to better-mem
		var relatedContext []MessageRelatedContext
		historyMessages := s.history.GetMessages(s.config.ChatID, 2)
		for _, msg := range historyMessages {
			relatedContext = append(relatedContext, MessageRelatedContext{
				User:    msg.Role,
				Context: msg.Content,
			})
		}

		if err := s.apiClient.SendMessage(s.config.ChatID, userInput, relatedContext); err != nil {
			fmt.Printf("Error sending message to better-mem: %v\n", err)
		}

		if len(memories) > 0 {
			fmt.Printf("\n%s\n", cyan("Memórias Relacionadas:"))
			sort.Slice(memories, func(i, j int) bool {
				return memories[i].Score > memories[j].Score
			})
			for _, mem := range memories {
				memoryTypeStr := "Curto Prazo"
				colorFunc := yellow
				if mem.MemoryType == 2 {
					memoryTypeStr = "Longo Prazo"
					colorFunc = orange
				}
				fmt.Printf("  %s [%s] (Score: %.2f): %s\n", 
					colorFunc("•"), 
					memoryTypeStr, 
					mem.Score, 
					mem.Text,
				)
				fmt.Printf("    %s: %s\n", cyan("Criada em"), mem.CreatedAt)
				for _, context := range mem.RelatedContext {
					fmt.Printf("    %s: %s\n", cyan("Usuário"), context.User)
					fmt.Printf("    %s: %s\n", cyan("Contexto"), context.Context)
				}
			}
			fmt.Println()
		}

		aiResponse, err := s.generateResponse(userInput, memories)
		if err != nil {
			fmt.Printf("Erro ao gerar resposta: %v\n", err)
			continue
		}

		fmt.Printf("%s %s\n\n", green("AI:"), aiResponse)

		s.history.AddMessage(s.config.ChatID, "user", userInput)
		s.history.AddMessage(s.config.ChatID, "assistant", aiResponse)

		if err := s.history.Save(s.historyPath); err != nil {
			fmt.Printf("Aviso: Erro ao salvar histórico: %v\n", err)
		}
	}

	return nil
}

func (s *ChatSession) generateResponse(userInput string, memories []ScoredMemory) (string, error) {
	systemPrompt := "You are a helpful assistant"
	
	if len(memories) > 0 {
		systemPrompt += fmt.Sprintf("\n\nThese are memories you have from past conversations with %s:\n", s.config.Name)
		sort.Slice(memories, func(i, j int) bool {
			return memories[i].Score > memories[j].Score
		})
		for _, mem := range memories {
			systemPrompt += fmt.Sprintf("- %s (relevance: %.2f, created at: %s)\n", mem.Text, mem.Score, mem.CreatedAt)
			systemPrompt += "  Related context:\n"
			for _, context := range mem.RelatedContext {
				systemPrompt += fmt.Sprintf("  - From: %s\n    %s\n", context.User, context.Context)
			}
		}
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	historyMessages := s.history.GetMessages(s.config.ChatID, s.config.ChatHistoryBuffer)
	for _, msg := range historyMessages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	resp, err := s.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    s.config.Model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta gerada")
	}

	return resp.Choices[0].Message.Content, nil
}

