package llm

type LLMProvider interface {
	GetCompletion(prompt string) (string, error)
	TestProvider() error
}
