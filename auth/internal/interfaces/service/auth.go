package service

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) error
	ValidateToken(ctx context.Context, token string) (string, string, error)
}
