package impl

import (
	"context"
	"errors"
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
	//check if active order exists (mb delete if on gateway checking too)
	currentOrder, err := s.db.GetCurrentOrder(ctx, order.UserID)
	if err != nil && err.Error() != "no active order found" {
		return err
	}

	if currentOrder != nil {
		return errors.New("active order already exists")
	}

	if err := s.db.CreateOrder(ctx, order); err != nil {
		return err
	}

	if err := s.broker.PublishCreatedOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *service) GetOrders(ctx context.Context, userID uuid.UUID) ([]*infra.Order, error) {
	//TODO: implement
	return nil, nil
}

func (s *service) GetOrderById(ctx context.Context, orderID uuid.UUID) (*infra.Order, error) {
	//TODO: implement
	return nil, nil
}

func (s *service) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	//TODO: implement
	return nil
}

func (s *service) FinishOrder(ctx context.Context, orderID uuid.UUID) error {
	//TODO: implement
	return nil
}
