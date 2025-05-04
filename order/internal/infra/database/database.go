package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order_service/internal/config"
	"order_service/internal/infra"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostgresDB struct {
	Logger *zap.Logger
	Db     *pgxpool.Pool
}

func New(logger *zap.Logger, cfg *config.Postgres) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("failed to create pgx pool", zap.Error(err))
		return nil, err
	}

	// Проверка соединения
	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to ping pgx pool", zap.Error(err))
		return nil, err
	}

	logger.Info("connected to database", zap.String("dsn", dsn))
	return &PostgresDB{Logger: logger, Db: pool}, nil
}

func (p *PostgresDB) Close() error {
	p.Db.Close()
	return nil
}

func (p *PostgresDB) GetCurrentOrder(ctx context.Context, userID uuid.UUID) (*infra.Order, error) {
	query := `
	SELECT
		order_id,
		user_id,
		order_address,
		order_location,
		order_date,
		order_time_gap,
		order_status
	FROM orders
	WHERE user_id = $1
	AND (order_status = 'pending' OR order_status = 'matching' OR order_status = 'signed')
	`
	var order infra.Order
	if err := p.Db.QueryRow(ctx, query, userID).Scan(&order.OrderID,
		&order.UserID,
		&order.OrderAddress,
		&order.OrderLocation,
		&order.OrderDate,
		&order.OrderTimeGap,
		&order.OrderStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no active order found")
		}
		return nil, fmt.Errorf("failed to get current order: %w", err)
	}
	return &order, nil
}

func (p *PostgresDB) CreateOrder(ctx context.Context, order *infra.Order) error {
	query := `
	INSERT INTO orders (
		user_id, 
		order_address, 
		order_location, 
		order_date, 
		order_time_gap, 
		order_status
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING order_id
	`

	err := p.Db.QueryRow(
		ctx,
		query,
		order.UserID,
		order.OrderAddress,
		order.OrderLocation,
		order.OrderDate,
		order.OrderTimeGap,
		order.OrderStatus,
	).Scan(&order.OrderID)

	if err != nil {
		p.Logger.Error("failed to create order", zap.Error(err))
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (p *PostgresDB) GetOrders(ctx context.Context, userID uuid.UUID) ([]*infra.Order, error) {
	query := `
	SELECT
		order_id,
		user_id,
		order_address,
		order_location,
		order_date,
		order_time_gap,
		order_status
	FROM orders
	WHERE user_id = $1
	`
	rows, err := p.Db.Query(ctx, query, userID)
	if err != nil {
		p.Logger.Error("failed to get orders", zap.Error(err))
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	orders := []*infra.Order{}
	for rows.Next() {
		var order infra.Order
		if err := rows.Scan(&order.OrderID,
			&order.UserID,
			&order.OrderAddress,
			&order.OrderLocation,
			&order.OrderDate,
			&order.OrderTimeGap,
			&order.OrderStatus,
		); err != nil {
			p.Logger.Error("failed to scan order", zap.Error(err))
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		p.Logger.Error("failed to iterate over orders", zap.Error(err))
		return nil, fmt.Errorf("failed to iterate over orders: %w", err)
	}

	return orders, nil
}

func (p *PostgresDB) GetOrderById(ctx context.Context, orderID uuid.UUID) (*infra.Order, error) {
	query := `
	SELECT
		order_id,
		user_id,
		order_address,
		order_location,
		order_date,
		order_time_gap,
		order_status
	FROM orders
	WHERE order_id = $1
	`
	var order infra.Order
	if err := p.Db.QueryRow(ctx, query, orderID).Scan(&order.OrderID,
		&order.UserID,
		&order.OrderAddress,
		&order.OrderLocation,
		&order.OrderDate,
		&order.OrderTimeGap,
		&order.OrderStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}

	return &order, nil
}
