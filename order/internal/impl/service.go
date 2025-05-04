package impl

import (
	"context"
	"order_service/internal/infra"
	"order_service/internal/infra/broker"
	"order_service/internal/infra/database"
	"order_service/internal/interfaces"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	db     *database.PostgresDB
	broker *broker.RabbitMQ
}

func New(logger *zap.Logger, db *database.PostgresDB, broker *broker.RabbitMQ) interfaces.Service {
	return &service{logger: logger, db: db, broker: broker}
}

func (s *service) CreateOrder(ctx context.Context, order *infra.Order) error {
	s.logger.Info("Creating order", zap.Any("order", order))

	// s.logger.Info("Getting current order")
	// currentOrder, err := s.db.GetCurrentOrder(ctx, order.UserID)
	// if err != nil && err.Error() != "no active order found" {
	// 	s.logger.Error("Failed to get current order", zap.Error(err))
	// 	return err
	// }

	// s.logger.Info("Checking if active order already exists")
	// if currentOrder != nil {
	// 	s.logger.Info("Active order already exists", zap.Any("order", currentOrder))
	// 	return errors.New("active order already exists")
	// }

	s.logger.Info("Creating order")
	if err := s.db.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return err
	}

	s.logger.Info("Publishing created order")
	if err := s.broker.PublishCreatedOrder(ctx, order); err != nil {
		s.logger.Error("Failed to publish created order", zap.Error(err))
		return err
	}

	s.logger.Info("Order created successfully")
	return nil
}

func (s *service) GetOrders(ctx context.Context, userID uuid.UUID) ([]*infra.Order, error) {
	s.logger.Info("Getting orders", zap.String("userID", userID.String()))
	orders, err := s.db.GetOrders(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get orders", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Orders found", zap.Any("orders", orders))

	return orders, nil
}

func (s *service) GetOrderById(ctx context.Context, orderID uuid.UUID) (*infra.Order, error) {
	s.logger.Info("Getting order by ID", zap.String("orderID", orderID.String()))
	order, err := s.db.GetOrderById(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to get order by ID", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Order found", zap.Any("order", order))

	return order, nil
}

func (s *service) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	//TODO: implement
	return nil
}

func (s *service) FinishOrder(ctx context.Context, orderID uuid.UUID) error {
	//TODO: implement
	return nil
}
