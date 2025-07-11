package interfaces

import (
	"context"
	"order_service/internal/infra"

	"github.com/google/uuid"
)

type Service interface {
	CreateOrder(ctx context.Context, order *infra.Order) error
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*infra.Order, error)
	GetAvailableOrders(ctx context.Context) ([]*infra.Order, error)
	GetOrderById(ctx context.Context, orderID uuid.UUID) (*infra.Order, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
	CompleteOrder(ctx context.Context, orderID uuid.UUID) error
}
