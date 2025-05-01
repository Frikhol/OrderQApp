package infra

import (
	"auth_service/internal/config"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type PostgresDB struct {
	Logger *zap.Logger
	Db     *sql.DB
}

func New(logger *zap.Logger, cfg *config.Config) (*PostgresDB, error) {
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

func (p *PostgresDB) UserExists(ctx context.Context, email string) error {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	var count int
	err := p.Db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		p.Logger.Error("database error checking if user exists", zap.Error(err))
		return fmt.Errorf("database error: %w", err)
	}
	if count == 0 {
		return errors.New("user does not exist")
	}
	return nil // No error means user exists
}

func (p *PostgresDB) InsertUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	_, err := p.Db.ExecContext(ctx, query, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	p.Logger.Info("getting user by email", zap.String("email", email))

	query := `SELECT email, password FROM users WHERE email = $1`
	var user User
	err := p.Db.QueryRowContext(ctx, query, email).Scan(&user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			p.Logger.Error("user not found", zap.String("email", email))
			return nil, errors.New("user not found")
		}
		p.Logger.Error("database error", zap.Error(err))
		return nil, fmt.Errorf("database error: %w", err)
	}

	p.Logger.Info("user found", zap.String("email", user.Email))
	return &user, nil
}
