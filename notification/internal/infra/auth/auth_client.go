package auth

import (
	"context"
	"errors"
	"notification_service/proto/auth_service"
)

type AuthClient struct {
	client auth_service.AuthServiceClient
}

func NewAuthClient(ac auth_service.AuthServiceClient) *AuthClient {
	return &AuthClient{ac}
}

func (ac *AuthClient) ValidateToken(token string) (string, error) {
	res, err := ac.client.ValidateToken(context.TODO(), &auth_service.ValidateTokenRequest{Token: token})
	if err != nil {
		return "", err
	}
	if res.Success == false {
		return "", errors.New("invalid or expired token")
	}
	return res.UserId, nil
}
