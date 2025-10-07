package config

type generalConfig struct {
	ApiPort          int
	InferenceAddress string
}

func newGeneralConfig() *generalConfig {
	apiPort := getInt("API_PORT", 8080)
	inferenceAddress := getString("INFERENCE_ADDRESS", "localhost:50051")

	return &generalConfig{
		ApiPort:          apiPort,
		InferenceAddress: inferenceAddress,
	}
}

var General = newGeneralConfig()
