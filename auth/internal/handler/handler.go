package handler

import (
	"auth_service/internal/interfaces/service"
	pb "auth_service/proto/auth_service"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewHandler(service service.AuthService) *Handler {
	return &Handler{service: service}
}

func (s *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (s *Handler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
	return &pb.RegisterResponse{Success: true}, nil
}

func (s *Handler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	user, role, err := s.service.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return &pb.ValidateTokenResponse{Success: true, UserId: user, Role: role}, nil
}
