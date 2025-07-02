package auth

type Auth interface {
	ValidateToken(token string) error
}
