//package main
//
//import (
//	"log"
//	"net"
//
//	"google.golang.org/grpc"
//	"orderq/"
//)
//
//func main() {
//	// cfg, err := config.GetConfigFromEnv()
//	// if err != nil {
//	// 	log.Fatalf("Failed to load configuration: %s\n", err.Error())
//	// }
//
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	s := grpc.NewServer()
//	pb.RegisterOrderServiceServer(s, &order.Server{})
//
//	log.Println("Starting Order Service on port 50051")
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}
