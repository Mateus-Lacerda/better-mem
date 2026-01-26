package main

import (
	docs "better-mem/docs"
	v1 "better-mem/internal/api/v1"
	"better-mem/internal/database/mongo"
	"better-mem/internal/database/qdrant"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const apiPort = "5042"

// @title Better Mem API
// @version 1.0
// @description This is the API for the Better Mem project.
// @contact.name Mateus Lacerda
// @contact.email mateuslacerda253@gmail.com
func main() {
	slog.Info("api", "message", "connecting to mongo")
	if err := mongo.TestMongo(); err != nil {
		slog.Error("failed to connect to mongo", "error", err)
		return
	}

	if err := qdrant.TestQdrant(); err != nil {
		slog.Error("failed to connect to qdrant", "error", err)
		return
	}

	startApi()
}

func startApi() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("api", "error", err)
			panic(err)
		}
	}()
	router := gin.Default()
	v1.Register(router)
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET(
		"/swagger/v1/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
	router.Run("0.0.0.0:" + apiPort)
}
