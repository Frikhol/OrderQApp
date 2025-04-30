package interfaces

import (
	"context"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) error
	ValidateToken(ctx context.Context, token string) error
}
