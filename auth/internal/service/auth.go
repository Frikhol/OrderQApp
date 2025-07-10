package service

import (
	infra2 "auth_service/internal/domain/models"
	"auth_service/internal/infra"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	logger *zap.Logger
	db     *infra.PostgresDB
	secret string
}

func NewService(logger *zap.Logger, db *infra.PostgresDB, secret string) *Service {
	return &Service{logger: logger, db: db, secret: secret}
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	//empty check
	if email == "" || password == "" {
		s.logger.Error("email or password is empty")
		return "", errors.New("email and password are required")
	}

	//email check
	if !strings.Contains(email, "@") {
		s.logger.Error("invalid email format", zap.String("email", email))
		return "", errors.New("invalid email")
	}

	s.logger.Info("attempting login", zap.String("email", email))

	//get user directly without checking existence first
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return "", errors.New("no such user")
	}

	s.logger.Info("user found, checking password", zap.String("email", user.Email), zap.String("password", user.Password))

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Error("password comparison failed", zap.Error(err))
		return "", errors.New("invalid password")
	}

	s.logger.Info("password correct, generating token")

	//create token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		s.logger.Error("failed to create token", zap.Error(err))
		return "", errors.New("failed to create token")
	}

	s.logger.Info("login successful", zap.String("email", email))
	return tokenString, nil
}

func (s *Service) Register(ctx context.Context, email string, password string) error {
	//empty check
	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	//email check
	if !strings.Contains(email, "@") {
		return errors.New("invalid email")
	}

	//user exists check
	err := s.db.UserExists(ctx, email)
	if err == nil {
		// No error means user exists
		return errors.New("user already exists")
	} else if !strings.Contains(err.Error(), "user does not exist") {
		return errors.New("error checking if user exists")
	}
	// If error is "user does not exist", proceed with registration

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	//create user
	user := infra2.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     infra2.ClientRole,
	}

	//save user
	err = s.db.InsertUser(ctx, &user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *Service) ValidateToken(ctx context.Context, tokenString string) (string, string, error) {
	s.logger.Info("validating token", zap.String("token", tokenString))
	//empty check
	if tokenString == "" {
		s.logger.Error("token is required")
		return "", "", errors.New("token is required")
	}

	//validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		s.logger.Error("invalid token", zap.Error(err))
		return "", "", errors.New("invalid token")
	}

	//check expiration
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || time.Now().Unix() > int64(claims["exp"].(float64)) {
		s.logger.Error("token expired")
		return "", "", errors.New("token expired")
	}

	s.logger.Info("token is valid", zap.String("user_id", claims["user_id"].(string)))

	return claims["user_id"].(string), claims["role"].(string), nil
}
