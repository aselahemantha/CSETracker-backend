package services

import "github.com/aselahemantha/cseTrackerBackend/auth-service/internal/models"

type AuthService interface {
	CheckEmail(email string) (bool, error)
	RequestVerification(email string) error
	VerifyCode(email, code string) (bool, error)
	CompleteRegistration(email, password, firstName, lastName, code string) (*models.User, error)
	Login(email, password string) (string, error)
	Register(email, password string) (*models.User, error)
}
