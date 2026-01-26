package llm

import (
	"fmt"
	"reflect"
)

// Returned when the llm provider is not running
func LLMProviderError(provider LLMProvider) error {
	providerName := reflect.TypeOf(provider).String()
	return fmt.Errorf("Provider %s is not running.", providerName)
}
