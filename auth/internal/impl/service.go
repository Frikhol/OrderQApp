package impl

import (
	"auth_service/internal/infra"
	"auth_service/internal/interfaces"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger *zap.Logger
	db     *infra.PostgresDB
	secret string
}

func New(logger *zap.Logger, db *infra.PostgresDB, secret string) interfaces.Service {
	return &service{logger: logger, db: db, secret: secret}
}

func (s *service) Login(ctx context.Context, email string, password string) (string, error) {
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

	//get client directly without checking existence first
	client, err := s.db.GetClientByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get client", zap.Error(err))
		return "", errors.New("no such client")
	}

	s.logger.Info("client found, checking password", zap.String("email", client.Email), zap.String("password", client.Password))

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(password)); err != nil {
		s.logger.Error("password comparison failed", zap.Error(err))
		return "", errors.New("invalid password")
	}

	s.logger.Info("password correct, generating token")

	//create token
	claims := jwt.MapClaims{
		"user_id": client.ID,
		"role":    client.Role,
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

func (s *service) AgentLogin(ctx context.Context, email string, password string) (string, error) {
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

	//get agent directly without checking existence first
	agent, err := s.db.GetAgentByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return "", errors.New("no such agent")
	}

	s.logger.Info("agent found, checking password", zap.String("email", agent.Email), zap.String("password", agent.Password))

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(agent.Password), []byte(password)); err != nil {
		s.logger.Error("password comparison failed", zap.Error(err))
		return "", errors.New("invalid password")
	}
	//create token
	claims := jwt.MapClaims{
		"user_id": agent.ID,
		"role":    agent.Role,
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

func (s *service) Register(ctx context.Context, email string, password string) error {
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
	user := infra.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     infra.ClientRole,
	}

	//save user
	err = s.db.InsertUser(ctx, &user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *service) ValidateToken(ctx context.Context, tokenString string) (string, string, error) {
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
