package config

type llmConfig struct {
	BaseUrl string
	Model   string
}

func newLlmConfig() llmConfig {
	baseUrl := getString("LLM_BASE_URL", "http://localhost:11434")
	model := getString("LLM_MODEL", "qwen2.5:7b")
	return llmConfig{BaseUrl: baseUrl, Model: model}
}

var Llm = newLlmConfig()
