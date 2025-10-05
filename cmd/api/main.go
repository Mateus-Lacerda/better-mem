package main

import (
	"log/slog"
	docs "better-mem/docs"
	v1 "better-mem/internal/api/v1"
	"better-mem/internal/database/mongo"
	"better-mem/internal/database/qdrant"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const apiPort = "5042"

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
	router.Run("0.0.0.0:"+apiPort)
}


func main() {
	slog.Info("api", "message", "connecting to mongo")
	if err := mongo.SetupMongo(); err != nil {
		slog.Error("failed to connect to mongo", "error", err)
		return
	}

	if err := qdrant.SetupQdrant(); err != nil {
		slog.Error("failed to connect to qdrant", "error", err)
		return
	}

	startApi()
}
