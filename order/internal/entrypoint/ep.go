package entrypoint

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"order_service/internal/config"
	handlers "order_service/internal/handlers"
	impl "order_service/internal/impl"
	"order_service/internal/infra/broker"
	"order_service/internal/infra/database"
	proto "order_service/proto/order_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	db, err := database.New(logger, &cfg.Postgres)
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
	}
	defer db.Close()

	broker, err := broker.New(logger, &cfg.RabbitMQ)
	if err != nil {
		logger.Fatal("failed to create broker", zap.Error(err))
	}
	defer broker.Close()

	grpcServer := grpc.NewServer()
	proto.RegisterOrderServiceServer(grpcServer, handlers.New(impl.New(logger, db, broker)))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Order service running", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed to serve", zap.Error(err))
		}
	}()

	<-done
	logger.Info("Order service stopped")
	grpcServer.Stop()

	return nil
}
