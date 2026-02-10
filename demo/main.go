package main

import (
	demo "github.com/Mateus-Lacerda/better-mem/demo/src"
	"fmt"
	"log"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

const (
	configPath      = "data/bettermem_config.json"
	chatHistoryPath = "data/chat_history.json"
)

func main() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	myFigure := figure.NewFigure("BetterMem", "slant", true)
	myFigure.Print()
	fmt.Println()

	config, err := demo.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("%s Erro ao carregar configuração: %v\n", red("✗"), err)
	}

	if config == nil {
		fmt.Println(cyan("=== Configuração Inicial ==="))
		config = demo.DefaultConfig()

		name, err := demo.PromptString("Qual é o seu nome?", "")
		if err != nil {
			log.Fatalf("%s Erro ao obter nome: %v\n", red("✗"), err)
		}
		config.Name = name
		config.ChatID = demo.CreateSlug(name)

		if err := demo.ConfigureProvider(config); err != nil {
			log.Fatalf("%s Erro ao configurar provedor: %v\n", red("✗"), err)
		}

		changeFetchConfig, err := demo.PromptYesNo("Deseja alterar as configurações padrão de busca de memórias?")
		if err != nil {
			log.Fatalf("%s Erro ao obter resposta: %v\n", red("✗"), err)
		}

		if changeFetchConfig {
			if err := demo.ConfigureSettings(config); err != nil {
				log.Fatalf("%s Erro ao configurar: %v\n", red("✗"), err)
			}
		}

		apiClient := demo.NewAPIClient(config.APIBaseURL)
		fmt.Printf("\n%s Criando chat no sistema...\n", cyan("→"))
		if err := apiClient.CreateChat(config.ChatID); err != nil {
			log.Fatalf("%s Erro ao criar chat: %v\n", red("✗"), err)
		}
		fmt.Printf("%s Chat criado com sucesso!\n", green("✓"))

		if err := config.Save(configPath); err != nil {
			log.Fatalf("%s Erro ao salvar configuração: %v\n", red("✗"), err)
		}
		fmt.Printf("%s Configuração salva!\n\n", green("✓"))
	}

	history, err := demo.LoadChatHistory(chatHistoryPath)
	if err != nil {
		log.Fatalf("%s Erro ao carregar histórico: %v\n", red("✗"), err)
	}

	apiClient := demo.NewAPIClient(config.APIBaseURL)

	for {
		choice, err := demo.ShowMainMenu()
		if err != nil {
			if err.Error() == "^C" {
				fmt.Println(cyan("\n\nAté logo!"))
				os.Exit(0)
			}
			log.Fatalf("%s Erro no menu: %v\n", red("✗"), err)
		}

		switch choice {
		case "Talk":
			chatSession := demo.NewChatSession(config, apiClient, history, chatHistoryPath)
			if err := chatSession.Start(); err != nil {
				fmt.Printf("%s Erro no chat: %v\n", red("✗"), err)
			}

		case "Configure":
			config.Print()
			if err := demo.ConfigureSettings(config); err != nil {
				fmt.Printf("%s Erro ao configurar: %v\n", red("✗"), err)
				continue
			}
			if err := config.Save(configPath); err != nil {
				fmt.Printf("%s Erro ao salvar configuração: %v\n", red("✗"), err)
				continue
			}
			fmt.Printf("%s Configuração atualizada!\n\n", green("✓"))

		case "Exit":
			fmt.Println(cyan("\nAté logo!"))
			os.Exit(0)
		}
	}
}
