package entrypoint

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/infra/auth"
	"notification_service/internal/infra/broker"
	"notification_service/internal/transport/websocket"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config, logger *zap.Logger) error {

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Notification service starting...")

	authClient := auth.NewAuthClient()
	rabbitClient := broker.NewRabbitClient()

	go func() {
		http.HandleFunc("/ws", websocket.NewHandler(authClient, rabbitClient).HandleUser)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Notification service running", zap.String("port", cfg.Port))
	}()

	<-done
	logger.Info("Notification service stopped")

	return nil
}
