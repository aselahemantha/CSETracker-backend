package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aselahemantha/cseTrackerBackend/auth-service/internal/handlers"
	"github.com/aselahemantha/cseTrackerBackend/auth-service/internal/repository"
	"github.com/aselahemantha/cseTrackerBackend/auth-service/internal/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("AUTH_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://auth_user:auth_password@localhost:5433/auth_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Initialize Repository
	userRepo := repository.NewPostgresUserRepository(db)
	if err := userRepo.InitTable(); err != nil {
		log.Fatalf("Error initializing users table: %v", err)
	}

	// Initialize Service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_super_secret_key"
	}
	authService := services.NewAuthService(userRepo, jwtSecret)

	// Initialize Handler
	authHandler := handlers.NewAuthHandler(authService)

	// Setup Router
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Auth Service is healthy")
	}).Methods("GET")

	r.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	log.Println("Starting Auth Service on :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
