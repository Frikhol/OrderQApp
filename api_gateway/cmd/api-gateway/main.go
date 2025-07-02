package main

import (
	"api_gateway/proto/auth_service"
	"api_gateway/proto/order_service"
	"api_gateway/router" // Import routers to initialize them
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/beego/beego/v2/server/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting server...")

		AuthConn, err := grpc.NewClient("auth_service:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Connected to auth service")

		defer func() {
			if err := AuthConn.Close(); err != nil {
				log.Println(err)
			}
		}()

		OrderConn, err := grpc.NewClient("order_service:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Connected to order service")

		defer func() {
			if err := OrderConn.Close(); err != nil {
				log.Println(err)
			}
		}()

		AuthClient := auth_service.NewAuthServiceClient(AuthConn)
		OrderClient := order_service.NewOrderServiceClient(OrderConn)

		router.InitRoutes(AuthClient, OrderClient)

		web.Run()
	}()

	// Wait for shutdown signal
	<-shutdownChan
	log.Println("Shutting down server...")

	log.Println("Server stopped")
}
