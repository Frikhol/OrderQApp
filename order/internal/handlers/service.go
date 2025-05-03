package handlers

import (
	"context"
	"order_service/internal/infra"
	"order_service/internal/interfaces"
	"order_service/internal/mapper"
	pb "order_service/proto/order_service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	service interfaces.Service
}

func New(service interfaces.Service) *OrderService {
	return &OrderService{service: service}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &infra.Order{
		UserID:        uuid.MustParse(req.GetUserId()),
		OrderAddress:  req.GetOrderAddress(),
		OrderLocation: req.GetOrderLocation(),
		OrderDate:     req.GetOrderDate().AsTime(),
		OrderTimeGap:  req.GetOrderTimeGap().AsDuration(),
		OrderStatus:   "pending",
	}

	err := s.service.CreateOrder(ctx, order)
	if err != nil {
		return &pb.CreateOrderResponse{Success: false}, status.Errorf(codes.Internal, "create order failed: %v", err)
	}

	return &pb.CreateOrderResponse{Success: true}, nil
}

func (s *OrderService) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	orders, err := s.service.GetOrders(ctx, uuid.MustParse(req.GetUserId()))
	if err != nil {
		return &pb.GetOrdersResponse{Orders: nil}, status.Errorf(codes.Internal, "get orders failed: %v", err)
	}

	return &pb.GetOrdersResponse{Orders: mapper.ToPbOrders(orders)}, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, req *pb.GetOrderByIdRequest) (*pb.GetOrderByIdResponse, error) {
	order, err := s.service.GetOrderById(ctx, uuid.MustParse(req.GetOrderId()))
	if err != nil {
		return &pb.GetOrderByIdResponse{Order: nil}, status.Errorf(codes.Internal, "get order by id failed: %v", err)
	}

	return &pb.GetOrderByIdResponse{Order: mapper.ToPbOrder(order)}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	err := s.service.CancelOrder(ctx, uuid.MustParse(req.GetOrderId()))
	if err != nil {
		return &pb.CancelOrderResponse{Success: false}, status.Errorf(codes.Internal, "cancel order failed: %v", err)
	}

	return &pb.CancelOrderResponse{Success: true}, nil
}

func (s *OrderService) FinishOrder(ctx context.Context, req *pb.FinishOrderRequest) (*pb.FinishOrderResponse, error) {
	err := s.service.FinishOrder(ctx, uuid.MustParse(req.GetOrderId()))
	if err != nil {
		return &pb.FinishOrderResponse{Success: false}, status.Errorf(codes.Internal, "finish order failed: %v", err)
	}

	return &pb.FinishOrderResponse{Success: true}, nil
}
