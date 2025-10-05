package core

type NewMessage struct {
	ChatId  string `json:"chat_id"`
	Message string `json:"message"`
}

type LabeledMessage struct {
	NewMessage       `json:"new_message"`
	Label            MemoryTypeEnum `json:"label"`
	MessageEmbedding []float32      `json:"message_embedding"`
}
