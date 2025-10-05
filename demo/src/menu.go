package demo

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func ShowMainMenu() (string, error) {
	prompt := promptui.Select{
		Label: "Menu Principal",
		Items: []string{"Talk", "Configure", "Exit"},
		Keys: &promptui.SelectKeys{
			Prev: promptui.Key{Code: promptui.KeyPrev, Display: "↑"},
			Next: promptui.Key{Code: promptui.KeyNext, Display: "↓"},
			PageUp: promptui.Key{Code: 'k', Display: "k"},
			PageDown: promptui.Key{Code: 'j', Display: "j"},
		},
	}

	_, result, err := prompt.Run()
	return result, err
}

func PromptString(label string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	return prompt.Run()
}

func PromptInt(label string, defaultValue int) (int, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: strconv.Itoa(defaultValue),
		Validate: func(input string) error {
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("deve ser um número inteiro")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result)
}

func PromptFloat(label string, defaultValue float32) (float32, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: fmt.Sprintf("%.2f", defaultValue),
		Validate: func(input string) error {
			_, err := strconv.ParseFloat(input, 32)
			if err != nil {
				return fmt.Errorf("deve ser um número decimal")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	val, _ := strconv.ParseFloat(result, 32)
	return float32(val), nil
}

func PromptYesNo(label string) (bool, error) {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Sim", "Não"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Sim", nil
}

func ConfigureSettings(config *Config) error {
	fmt.Println("\n=== Configurações ===")
	
	change, err := PromptYesNo("Deseja alterar o Limite de Memórias?")
	if err != nil {
		return err
	}
	if change {
		config.Limit, err = PromptInt("Limite de Memórias", config.Limit)
		if err != nil {
			return err
		}
	}

	change, err = PromptYesNo("Deseja alterar o Limite de Busca Vetorial?")
	if err != nil {
		return err
	}
	if change {
		config.VectorSearchLimit, err = PromptInt("Limite de Busca Vetorial", config.VectorSearchLimit)
		if err != nil {
			return err
		}
	}

	change, err = PromptYesNo("Deseja alterar o Threshold de Busca Vetorial?")
	if err != nil {
		return err
	}
	if change {
		config.VectorSearchThreshold, err = PromptFloat("Threshold de Busca Vetorial", config.VectorSearchThreshold)
		if err != nil {
			return err
		}
	}

	change, err = PromptYesNo("Deseja alterar o Threshold de Longo Prazo?")
	if err != nil {
		return err
	}
	if change {
		config.LongTermThreshold, err = PromptFloat("Threshold de Longo Prazo", config.LongTermThreshold)
		if err != nil {
			return err
		}
	}

	change, err = PromptYesNo("Deseja alterar o Buffer de Histórico?")
	if err != nil {
		return err
	}
	if change {
		config.ChatHistoryBuffer, err = PromptInt("Buffer de Histórico", config.ChatHistoryBuffer)
		if err != nil {
			return err
		}
	}

	change, err = PromptYesNo("Deseja alterar o Modelo?")
	if err != nil {
		return err
	}
	if change {
		config.Model, err = PromptString("Modelo", config.Model)
		if err != nil {
			return err
		}
	}

	change, err = PromptYesNo("Deseja alterar a URL da API?")
	if err != nil {
		return err
	}
	if change {
		config.APIBaseURL, err = PromptString("URL da API", config.APIBaseURL)
		if err != nil {
			return err
		}
	}

	return nil
}

