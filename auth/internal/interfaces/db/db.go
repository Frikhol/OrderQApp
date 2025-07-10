package db

import (
	"auth_service/internal/domain/models"
	"github.com/google/uuid"
)

type Database interface {
	UserExists(username uuid.UUID) (bool, error)
	GetUser(username uuid.UUID) (models.User, error)
	InsertUser(user models.User) (bool, error)
}
