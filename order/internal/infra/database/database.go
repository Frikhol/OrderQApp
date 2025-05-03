package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order_service/internal/config"
	"order_service/internal/infra"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostgresDB struct {
	Logger *zap.Logger
	Db     *sql.DB
}

func New(logger *zap.Logger, cfg *config.Postgres) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("connected to database", zap.String("dsn", dsn))
	return &PostgresDB{Logger: logger, Db: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.Db.Close()
}

func (p *PostgresDB) GetCurrentOrder(ctx context.Context, userID uuid.UUID) (*infra.Order, error) {
	query := `
	SELECT * FROM orders WHERE user_id = $1 AND (order_status = 'pending' OR order_status = 'matching' OR order_status = 'signed')
	`
	var order infra.Order
	if err := p.Db.QueryRowContext(ctx, query, userID).Scan(&order); err != nil {
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

	err := p.Db.QueryRowContext(
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
