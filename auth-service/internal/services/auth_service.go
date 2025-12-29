package services

import (
	"errors"
	"time"

	"github.com/aselahemantha/cseTrackerBackend/auth-service/internal/models"
	"github.com/aselahemantha/cseTrackerBackend/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo      repository.UserRepository
	JWTSecret []byte
}

func NewAuthService(repo repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		Repo:      repo,
		JWTSecret: []byte(secret),
	}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	existingUser, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
