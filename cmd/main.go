package main

import (
	"cart-api/internal/app"
	"cart-api/pkg/logger"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	if err := app.Run(logger); err != nil {
		logger.Fatal("Error starting app", zap.Error(err))
	}
}
