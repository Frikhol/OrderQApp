package infra

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID       uuid.UUID     `json:"order_id"`
	UserID        uuid.UUID     `json:"user_id"`
	OrderAddress  string        `json:"order_address"`
	OrderLocation string        `json:"order_location"`
	OrderDate     time.Time     `json:"order_date"`
	OrderTimeGap  time.Duration `json:"order_time_gap"`
	OrderStatus   string        `json:"order_status"`
}
