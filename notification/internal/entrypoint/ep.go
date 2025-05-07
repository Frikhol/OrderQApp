package entrypoint

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"notification_service/internal/config"
	handlers "notification_service/internal/handlers"
	impl "notification_service/internal/impl"
	"notification_service/internal/infra/broker"
	proto "notification_service/proto/notification_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	broker, err := broker.New(logger, &cfg.RabbitMQ)
	if err != nil {
		logger.Fatal("failed to create broker", zap.Error(err))
	}
	defer broker.Close()

	grpcServer := grpc.NewServer()
	proto.RegisterNotificationServiceServer(grpcServer, handlers.New(impl.New(logger, broker)))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			logger.Fatal("failed to listen", zap.Error(err))
		}
		logger.Info("Notification service started", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed to serve", zap.Error(err))
		}
	}()

	go func() {
		err := impl.New(logger, broker).HandleMessages()
		if err != nil {
			logger.Error("failed to handle messages", zap.Error(err))
		}
	}()

	<-done
	logger.Info("Notification service stopped")
	grpcServer.Stop()

	return nil
}
