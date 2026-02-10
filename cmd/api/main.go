package main

import (
	docs "github.com/Mateus-Lacerda/better-mem/docs"
	v1 "github.com/Mateus-Lacerda/better-mem/internal/api/v1"
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
	if err := setup(); err != nil {
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
