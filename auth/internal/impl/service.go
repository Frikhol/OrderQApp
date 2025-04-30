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
		return "", errors.New("email and password are required")
	}

	//email check
	if !strings.Contains(email, "@") {
		return "", errors.New("invalid email")
	}

	//get user
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("failed to get user")
	}

	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	//create token
	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", errors.New("failed to create token")
	}

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
	if s.db.UserExists(ctx, email) != nil {
		return errors.New("user already exists")
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	//create user
	user := infra.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	//save user
	err = s.db.InsertUser(ctx, &user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *service) ValidateToken(ctx context.Context, tokenString string) error {
	//empty check
	if tokenString == "" {
		return errors.New("token is required")
	}

	//validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	//check expiration
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || time.Now().Unix() > int64(claims["exp"].(float64)) {
		return errors.New("token expired")
	}

	return nil
}
