package core

// Additional textual context that might be related to the memory
type MessageRelatedContext struct {
	// The text that was used to generate the memory
	Context string `json:"context"`
	// User that generated the context, might be a name
	// or simply "user and assistant"
	User string `json:"user"`
}

// A new message
type NewMessage struct {
	// The chat id that the message belongs to
	ChatId string `json:"chat_id"`
	// The message text
	Message string `json:"message"`
	// The related context
	RelatedContext []MessageRelatedContext `json:"related_context"`
}

type LabeledMessage struct {
	NewMessage       `json:"new_message"`
	Label            MemoryTypeEnum `json:"label"`
	MessageEmbedding []float32      `json:"message_embedding"`
}
