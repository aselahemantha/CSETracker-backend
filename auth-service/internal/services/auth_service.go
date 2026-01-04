package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
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

// CheckEmail returns true if email exists, false otherwise
func (s *AuthService) CheckEmail(email string) (bool, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (s *AuthService) RequestVerification(email string) error {
	// Generate 6-digit code
	code, err := generateRandomCode()
	if err != nil {
		return err
	}

	// Store code
	if err := s.Repo.StoreVerificationCode(email, code); err != nil {
		return err
	}

	// Mock send email
	log.Printf("Sending verification code %s to %s", code, email)
	return nil
}

func (s *AuthService) VerifyCode(email, code string) (bool, error) {
	storedCode, err := s.Repo.GetVerificationCode(email)
	if err != nil {
		return false, err
	}
	if storedCode == "" {
		return false, errors.New("verification code expired or invalid")
	}
	return storedCode == code, nil
}

func (s *AuthService) CompleteRegistration(email, password, firstName, lastName, code string) (*models.User, error) {
	// Verify code again to be sure
	valid, err := s.VerifyCode(email, code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid verification code")
	}

	// Check if user already exists
	exists, err := s.CheckEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: string(hashedPassword),
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Clean up code
	_ = s.Repo.DeleteVerificationCode(email)

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

// Deprecated: Use CompleteRegistration
func (s *AuthService) Register(email, password string) (*models.User, error) {
	return s.CompleteRegistration(email, password, "", "", "mock_code") // This needs refactor if keeping backward compatibility, but for this task we switch flow
}

func generateRandomCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
