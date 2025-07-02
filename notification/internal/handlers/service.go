package handlers

import (
	"context"
	"notification_service/internal/interfaces"
	pb "notification_service/proto/notification_service"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	service interfaces.Service
}

func New(service interfaces.Service) *NotificationService {
	return &NotificationService{service: service}
}

func (s *NotificationService) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: "OK"}, nil
}
