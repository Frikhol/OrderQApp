package main

import (
	"log"

	"order_service/internal/config"
	"order_service/internal/entrypoint"
	"order_service/internal/logger"

	_ "github.com/lib/pq" // PostgreSQL driver
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatalf("config.GetConfigFromEnv: %s\n", err.Error())
	}

	zapLogger := logger.NewClientZapLogger(cfg.LogLevel, cfg.ServiceName)

	if err = entrypoint.Run(cfg, zapLogger); err != nil {
		zapLogger.Fatal("entrypoint.Run: ", zap.Error(err))
	}
}
