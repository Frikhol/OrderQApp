package main

import (
	"api_gateway/proto/auth_service"
	_ "api_gateway/routers" // Import routers to initialize them
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/beego/beego/v2/server/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	AuthConn   *grpc.ClientConn
	AuthClient auth_service.AuthServiceClient
)

func main() {
	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting server...")

		AuthConn, err := grpc.NewClient("auth:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := AuthConn.Close(); err != nil {
				log.Println(err)
			}
		}()

		AuthClient = auth_service.NewAuthServiceClient(AuthConn)

		web.Run()
	}()

	// Wait for shutdown signal
	<-shutdownChan
	log.Println("Shutting down server...")

	log.Println("Server stopped")
}
