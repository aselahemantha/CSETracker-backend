package repository

import "github.com/aselahemantha/cseTrackerBackend/auth-service/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	InitTable() error
	StoreVerificationCode(email, code string) error
	GetVerificationCode(email string) (string, error)
	DeleteVerificationCode(email string) error
}
