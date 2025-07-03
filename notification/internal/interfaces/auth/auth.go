package auth

type Auth interface {
	ValidateToken(token string) (string, error) //returns userId
}
