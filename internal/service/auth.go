package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Vladimirmoscow84/Event_Booker/internal/model"
	"github.com/Vladimirmoscow84/Event_Booker/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage   *storage.Storage
	jwtSecret string
}

// NewAuthService конструктор AuthService
func NewAuthService(storage *storage.Storage, jwtSecret string) *AuthService {
	return &AuthService{
		storage:   storage,
		jwtSecret: jwtSecret,
	}
}

// hashPassword генерирует bcrypt-хеш пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[service-auth] failed to generate hash password: %w", err)
	}
	return string(bytes), nil
}

// Register создаёт нового пользователя с хешированным паролем
func (a *AuthService) Register(ctx context.Context, email, password string) (int, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("[service-auth] failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
	}

	id, err := a.storage.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("[service-auth] failed to register user: %w", err)
	}

	return id, nil
}

// Login проверяет email и пароль, возвращает JWT
func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.storage.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("[service-auth] user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("[service-auth] invalid credentials")
	}

	token, err := a.generateJWT(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("auth: failed to generate JWT: %w", err)
	}

	return token, nil
}

// generateJWT создаёт JWT для пользователя
func (a *AuthService) generateJWT(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}
