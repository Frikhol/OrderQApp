package infra

import (
	"auth_service/internal/config"
	"auth_service/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"

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

func (p *PostgresDB) UserExists(email string) error {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	var count int
	err := p.Db.QueryRowContext(context.TODO(), query, email).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return errors.New("database error" + err.Error())
	}
	if count == 0 {
		return errors.New("user does not exist")
	}
	return nil // No error means user exists
}

func (p *PostgresDB) InsertUser(user *models.User) error {
	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`
	_, err := p.Db.ExecContext(context.TODO(), query, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) GetUserById(id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, password, role FROM users WHERE id = $1`
	var user models.User
	err := p.Db.QueryRowContext(context.TODO(), query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}
	return &user, nil
}

func (p *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	var user models.User
	err := p.Db.QueryRowContext(context.TODO(), query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}

	return &user, nil
}
