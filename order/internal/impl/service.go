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

	s.logger.Info("Creating order")
	if err := s.db.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return err
	}

	s.logger.Info("Publishing created order")
	if err := s.broker.PublishOrderCreated(ctx, order); err != nil {
		s.logger.Error("Failed to publish created order", zap.Error(err))
		return err
	}

	s.logger.Info("Order created successfully")
	return nil
}

func (s *service) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*infra.Order, error) {
	s.logger.Info("Getting orders", zap.String("userID", userID.String()))
	orders, err := s.db.GetUserOrders(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get orders", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Orders found", zap.Any("orders", orders))

	return orders, nil
}

func (s *service) GetAvailableOrders(ctx context.Context) ([]*infra.Order, error) {
	s.logger.Info("Getting available orders")
	orders, err := s.db.GetAvailableOrders(ctx)
	if err != nil {
		s.logger.Error("Failed to get available orders", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Available orders found", zap.Any("orders", orders))

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
	s.logger.Info("Cancelling order", zap.String("orderID", orderID.String()))

	if err := s.db.CancelOrder(ctx, orderID); err != nil {
		s.logger.Error("Failed to cancel order", zap.Error(err))
		return err
	}

	s.logger.Info("Publishing cancelled order")
	order, err := s.db.GetOrderById(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to get order by ID", zap.Error(err))
		return err
	}

	if err := s.broker.PublishOrderCancelled(ctx, order); err != nil {
		s.logger.Error("Failed to publish cancelled order", zap.Error(err))
		return err
	}

	s.logger.Info("Order cancelled successfully")
	return nil
}

func (s *service) CompleteOrder(ctx context.Context, orderID uuid.UUID) error {
	s.logger.Info("Completing order", zap.String("orderID", orderID.String()))

	if err := s.db.CompleteOrder(ctx, orderID); err != nil {
		s.logger.Error("Failed to complete order", zap.Error(err))
		return err
	}

	s.logger.Info("Publishing completed order")
	order, err := s.db.GetOrderById(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to get order by ID", zap.Error(err))
		return err
	}

	if err := s.broker.PublishOrderCompleted(ctx, order); err != nil {
		s.logger.Error("Failed to publish completed order", zap.Error(err))
		return err
	}

	s.logger.Info("Order finished successfully")
	return nil
}
