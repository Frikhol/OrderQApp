package infra

import "github.com/google/uuid"

type Role string

const (
	ClientRole  Role = "client"
	AgentRole   Role = "agent"
	ManagerRole Role = "manager"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
}
