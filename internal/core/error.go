package core

import (
	"errors"
)

var (
	// Returned when an operation references an unexistent chat
	ChatNotFound = errors.New("Chat not found")
	// Returned when trying to create a chat with an existing external id
	ChatExternalIdAlreadyExists = errors.New("This external id is already taken")
	// Returned when there is an unexpected error on message classification
	UnexpectedClassificationError = errors.New("Unexpected Classification Error")
)
