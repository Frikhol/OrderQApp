package auth

type AuthClient struct {
}

func NewAuthClient() *AuthClient {
	return &AuthClient{}
}

func (ac *AuthClient) ValidateToken(token string) (string, error) {
	return "", nil
}
