package models

import "github.com/google/uuid"

type Role string

const (
	ClientRole Role = "client"
	AgentRole  Role = "agent"
	AdminRole  Role = "admin"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
}
