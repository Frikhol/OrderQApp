package database

import (
	"auth_service/internal/domain/models"
	"github.com/google/uuid"
)

type Database interface {
	UserExists(email string) error
	GetUserById(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(user *models.User) error
}
