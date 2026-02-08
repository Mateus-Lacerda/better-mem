//go:build server

package main

import (
	"better-mem/internal/database/mongo"
	"better-mem/internal/database/qdrant"
	"log/slog"
)

func setup() {
	slog.Info("api", "message", "connecting to mongo")
	if err := mongo.TestMongo(); err != nil {
		slog.Error("failed to connect to mongo", "error", err)
		return err
	}

	if err := qdrant.TestQdrant(); err != nil {
		slog.Error("failed to connect to qdrant", "error", err)
		return err
	}
}
