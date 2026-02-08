package demo

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatSession struct {
	config      *Config
	apiClient   *APIClient
	openaiClient *openai.Client
	history     *ChatHistory
	historyPath string
}

func NewChatSession(
	config *Config,
	apiClient *APIClient,
	history *ChatHistory,
	historyPath string,
) *ChatSession {
	var clientOptions []option.RequestOption

	if config.Provider == "ollama" {
		clientOptions = append(clientOptions, option.WithBaseURL(config.OllamaURL+"/v1"))
		clientOptions = append(clientOptions, option.WithAPIKey("ollama"))
	} else {
		clientOptions = append(clientOptions, option.WithAPIKey(config.OpenAIKey))
		clientOptions = append(clientOptions, option.WithBaseURL(config.APIBaseURL))
	}

	client := openai.NewClient(clientOptions...)

	return &ChatSession{
		config:       config,
		apiClient:    apiClient,
		openaiClient: &client,
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

		var memories []ScoredMemory
		err = s.withBetterMemLoading(func() error {
			var errFetch error
			memories, errFetch = s.apiClient.FetchMemories(s.config.ChatID, MemoryFetchRequest{
				Text:                  userInput,
				Limit:                 s.config.Limit,
				VectorSearchLimit:     s.config.VectorSearchLimit,
				VectorSearchThreshold: s.config.VectorSearchThreshold,
				LongTermThreshold:     s.config.LongTermThreshold,
			})
			if errFetch != nil {
				return errFetch
			}

			var relatedContext []MessageRelatedContext
			historyMessages := s.history.GetMessages(s.config.ChatID, 2)
			for _, msg := range historyMessages {
				relatedContext = append(relatedContext, MessageRelatedContext{
					User:    msg.Role,
					Context: msg.Content,
				})
			}

			return s.apiClient.SendMessage(s.config.ChatID, userInput, relatedContext)
		})

		if err != nil {
			fmt.Printf("%s %v\n", color.New(color.FgRed).Sprint("Error better-mem:"), err)
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

		aiResponse, err := s.withLLMLoading(func(yield func(string)) (string, error) {
			return s.generateResponseStreaming(userInput, memories, yield)
		})
		if err != nil {
			fmt.Printf("Erro ao gerar resposta: %v\n", err)
			continue
		}

		s.history.AddMessage(s.config.ChatID, "user", userInput)
		s.history.AddMessage(s.config.ChatID, "assistant", aiResponse)

		if err := s.history.Save(s.historyPath); err != nil {
			fmt.Printf("Aviso: Erro ao salvar histórico: %v\n", err)
		}
	}

	return nil
}

func (s *ChatSession) withBetterMemLoading(action func() error) error {
	start := time.Now()
	done := make(chan bool)
	
	cyan := color.New(color.FgCyan).SprintFunc()
	
	frames := []string{
		"( ●    )",
		"(  ●   )",
		"(   ●  )",
		"(    ● )",
		"(     ●)",
		"(    ● )",
		"(   ●  )",
		"(  ●   )",
		"( ●    )",
		"(●     )",
	}

	go func() {
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s %s", cyan(frames[i%len(frames)]), cyan("better-mem is searching..."))
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	err := action()
	done <- true
	duration := time.Since(start)
	
	// Clear line and print final result
	fmt.Printf("\r\033[2K\r")
	fmt.Printf("%.2fs %s\n", duration.Seconds(), "􀯐")
	
	return err
}

func (s *ChatSession) withLLMLoading(action func(func(string)) (string, error)) (string, error) {
	start := time.Now()
	done := make(chan bool)
	
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%.1fs ", time.Since(start).Seconds())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	green := color.New(color.FgGreen).SprintFunc()
	var fullResponse strings.Builder
	var firstChunk bool = true

	res, err := action(func(chunk string) {
		if firstChunk {
			done <- true
			duration := time.Since(start)
			fmt.Printf("\r\033[2K\r")
			fmt.Printf("%.2fs %s\n", duration.Seconds(), "􀌥")
			fmt.Printf("%s ", green("AI:"))
			firstChunk = false
		}
		fmt.Print(chunk)
		fullResponse.WriteString(chunk)
	})

	if err != nil {
		if firstChunk {
			done <- true
		}
		return "", err
	}

	duration := time.Since(start)
	fmt.Println()

	// Simple token estimation (words * 1.3 as a rule of thumb for English/Code)
	words := len(strings.Fields(fullResponse.String()))
	tokenCount := int(float64(words) * 1.3)
	if tokenCount == 0 && fullResponse.Len() > 0 {
		tokenCount = 1
	}
	
	tps := float64(tokenCount) / duration.Seconds()
	fmt.Printf("%s %.2f t/s\n\n", color.New(color.FgHiBlack).Sprint("TPS:"), tps)

	return res, nil
}

func (s *ChatSession) generateResponseStreaming(userInput string, memories []ScoredMemory, yield func(string)) (string, error) {
	var systemPrompt strings.Builder;

	systemPrompt.WriteString("You are a helpful assistant")
	
	if len(memories) > 0 {
		fmt.Fprintf(&systemPrompt, "\n\nThese are memories you have from past conversations with %s:\n", s.config.Name)
		sort.Slice(memories, func(i, j int) bool {
			return memories[i].Score > memories[j].Score
		})
		for _, mem := range memories {
			fmt.Fprintf(&systemPrompt, "- %s (relevance: %.2f, created at: %s)\n", mem.Text, mem.Score, mem.CreatedAt)
			systemPrompt .WriteString("  Related context:\n")
			for _, context := range mem.RelatedContext {
				fmt.Fprintf(&systemPrompt, "  - From: %s\n    %s\n", context.User, context.Context)
			}
		}
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt.String()),
	}

	historyMessages := s.history.GetMessages(s.config.ChatID, s.config.ChatHistoryBuffer)
	for _, msg := range historyMessages {
		switch msg.Role {
		case "user":
			messages = append(messages, openai.UserMessage(
				msg.Content,
			))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(
				msg.Content,
			))
		}
	}

	messages = append(messages, openai.UserMessage(
		userInput,
	))

	stream := s.openaiClient.Chat.Completions.NewStreaming(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:    s.config.Model,
			Messages: messages,
		},
	)

	var fullResponse strings.Builder
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			yield(content)
			fullResponse.WriteString(content)
		}
	}

	if err := stream.Err(); err != nil {
		return "", err
	}

	return fullResponse.String(), nil
}

func (s *ChatSession) generateResponse(userInput string, memories []ScoredMemory) (string, error) {
	var systemPrompt strings.Builder;

	systemPrompt.WriteString("You are a helpful assistant")
	
	if len(memories) > 0 {
		fmt.Fprintf(&systemPrompt, "\n\nThese are memories you have from past conversations with %s:\n", s.config.Name)
		sort.Slice(memories, func(i, j int) bool {
			return memories[i].Score > memories[j].Score
		})
		for _, mem := range memories {
			fmt.Fprintf(&systemPrompt, "- %s (relevance: %.2f, created at: %s)\n", mem.Text, mem.Score, mem.CreatedAt)
			systemPrompt .WriteString("  Related context:\n")
			for _, context := range mem.RelatedContext {
				fmt.Fprintf(&systemPrompt, "  - From: %s\n    %s\n", context.User, context.Context)
			}
		}
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt.String()),
	}

	historyMessages := s.history.GetMessages(s.config.ChatID, s.config.ChatHistoryBuffer)
	for _, msg := range historyMessages {
		switch msg.Role {
		case "user":
			messages = append(messages, openai.UserMessage(
				msg.Content,
			))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(
				msg.Content,
			))
		}
	}

	messages = append(messages, openai.UserMessage(
		userInput,
	))

	resp, err := s.openaiClient.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
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

