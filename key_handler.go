package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/m4trixdev/keygen-service/config"
	"github.com/m4trixdev/keygen-service/internal/models"
	"github.com/m4trixdev/keygen-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(username, password, role string) (*models.User, error) {
	if s.repo.UsernameExists(username) {
		return nil, errors.New("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashed),
		Role:     role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID.String(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(config.C.JWTSecret))
}
