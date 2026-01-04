package repository

import (
	"database/sql"
	"fmt"

	"github.com/aselahemantha/cseTrackerBackend/auth-service/internal/models"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) InitTable() error {
	userQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := r.DB.Exec(userQuery); err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	verificationQuery := `
	CREATE TABLE IF NOT EXISTS verification_codes (
		email VARCHAR(255) PRIMARY KEY,
		code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`
	if _, err := r.DB.Exec(verificationQuery); err != nil {
		return fmt.Errorf("error creating verification_codes table: %v", err)
	}

	return nil
}

func (r *PostgresUserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, first_name, last_name, password_hash) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.DB.QueryRow(query, user.Email, user.FirstName, user.LastName, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, first_name, last_name, password_hash, created_at, updated_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	user := &models.User{}
	var firstName, lastName sql.NullString
	err := row.Scan(&user.ID, &user.Email, &firstName, &lastName, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("error getting user by email: %v", err)
	}
	if firstName.Valid {
		user.FirstName = firstName.String
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}
	return user, nil
}

func (r *PostgresUserRepository) StoreVerificationCode(email, code string) error {
	query := `
	INSERT INTO verification_codes (email, code, expires_at) 
	VALUES ($1, $2, NOW() + INTERVAL '10 minutes')
	ON CONFLICT (email) DO UPDATE 
	SET code = $2, expires_at = NOW() + INTERVAL '10 minutes';`
	_, err := r.DB.Exec(query, email, code)
	if err != nil {
		return fmt.Errorf("error storing verification code: %v", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetVerificationCode(email string) (string, error) {
	query := `SELECT code FROM verification_codes WHERE email = $1 AND expires_at > NOW()`
	var code string
	err := r.DB.QueryRow(query, email).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("error getting verification code: %v", err)
	}
	return code, nil
}

func (r *PostgresUserRepository) DeleteVerificationCode(email string) error {
	query := `DELETE FROM verification_codes WHERE email = $1`
	_, err := r.DB.Exec(query, email)
	if err != nil {
		return fmt.Errorf("error deleting verification code: %v", err)
	}
	return nil
}
