package entrypoint

import (
	impl "auth_service/internal/service/service_impl"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"auth_service/internal/config"
	handlers "auth_service/internal/handlers"
	"auth_service/internal/infra"
	proto "auth_service/proto/auth_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	db, err := infra.New(logger, cfg)
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
	}
	defer db.Close()

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, handlers.New(impl.New(logger, db, cfg.Secret)))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Auth service running", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed to serve", zap.Error(err))
		}
	}()

	<-done
	logger.Info("Auth service stopped")
	grpcServer.Stop()

	return nil
}
