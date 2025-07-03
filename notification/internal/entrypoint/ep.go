package entrypoint

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"notification_service/internal/config"
	"notification_service/internal/connstore"
	"notification_service/internal/infra/auth"
	"notification_service/internal/infra/broker/rabbit"
	"notification_service/internal/transport/websocket"
	"notification_service/proto/auth_service"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config, logger *zap.Logger) error {

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Notification service starting...")

	AuthConn, err := grpc.NewClient("auth_service:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("grpc.NewClient: ", zap.Error(err))
	}

	authClient := auth.NewAuthClient(auth_service.NewAuthServiceClient(AuthConn))
	rabbitClient := rabbit.NewRabbitClient(cfg.RabbitMQ)

	err = rabbitClient.StartConsuming(context.TODO())
	if err != nil {
		logger.Fatal("rabbitClient.StartConsuming: ", zap.Error(err))
	}

	store := connstore.New()

	go func() {
		http.HandleFunc("/ws", websocket.NewHandler(authClient, store).HandleUser)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Notification service running", zap.String("port", cfg.Port))
	}()

	<-done
	logger.Info("Notification service stopped")

	return nil
}
