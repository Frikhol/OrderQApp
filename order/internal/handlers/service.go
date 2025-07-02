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

func (s *OrderService) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest) (*pb.GetUserOrdersResponse, error) {
	orders, err := s.service.GetUserOrders(ctx, uuid.MustParse(req.GetUserId()))
	if err != nil {
		return &pb.GetUserOrdersResponse{Orders: nil}, status.Errorf(codes.Internal, "get orders failed: %v", err)
	}

	return &pb.GetUserOrdersResponse{Orders: mapper.ToPbOrders(orders)}, nil
}

func (s *OrderService) GetAvailableOrders(ctx context.Context, req *pb.GetAvailableOrdersRequest) (*pb.GetAvailableOrdersResponse, error) {
	orders, err := s.service.GetAvailableOrders(ctx)
	if err != nil {
		return &pb.GetAvailableOrdersResponse{Orders: nil}, status.Errorf(codes.Internal, "get orders failed: %v", err)
	}

	return &pb.GetAvailableOrdersResponse{Orders: mapper.ToPbOrders(orders)}, nil
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

func (s *OrderService) CompleteOrder(ctx context.Context, req *pb.CompleteOrderRequest) (*pb.CompleteOrderResponse, error) {
	err := s.service.CompleteOrder(ctx, uuid.MustParse(req.GetOrderId()))
	if err != nil {
		return &pb.CompleteOrderResponse{Success: false}, status.Errorf(codes.Internal, "complete order failed: %v", err)
	}

	return &pb.CompleteOrderResponse{Success: true}, nil
}
