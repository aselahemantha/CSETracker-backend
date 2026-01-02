# CSE Tracker Backend

This repository contains the backend microservices for the CSE Tracker application. It is built using Go and follows a microservices architecture.

## Architecture

The system consists of the following services:

| Service | Port | Description | Dependencies |
|---------|------|-------------|--------------|
| **Auth Service** | `:8081` | User authentication and registration | Postgres (auth-db) |
| **Portfolio Service** | `:8082` | User portfolio management | Postgres (portfolio-db) |
| **Market Service** | `:8083` | Market data processing and caching | Redis (market-redis) |
| **API Gateway** | TBD | Unified entry point (Under Development) | - |

## Prerequisites

*   **Go**: Version 1.25 or later
*   **Docker & Docker Compose**: For running the database and cache containers.

## Getting Started

### 1. Start Infrastructure

Use Docker Compose to start the necessary databases and Redis cache.

```bash
docker-compose up -d
```

This will spin up:
*   `auth-db`: PostgreSQL on port `5433`
*   `portfolio-db`: PostgreSQL on port `5434`
*   `market-redis`: Redis on port `6379`

### 2. Run Microservices

You can run each service individually using Go. It is recommended to run them in separate terminal windows.

#### Auth Service

```bash
cd auth-service
go run cmd/api/main.go
```
*   **Health Check**: `http://localhost:8081/health`
*   **Endpoints**:
    *   `POST /auth/register`
    *   `POST /auth/login`

#### Portfolio Service

```bash
cd portfolio-service
go run cmd/api/main.go
```
*   **Health Check**: `http://localhost:8082/health`

#### Market Service

```bash
cd market-service
go run cmd/api/main.go
```
*   **Health Check**: `http://localhost:8083/health`

### 3. Workspace

This project uses a Go Workspace (`go.work`) to manage multi-module development. You can run commands from the root if configured, but running from each service directory is standard.

## Configuration

Services come with default configurations for local development. You can override them using environment variables.

### Auth Service Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_DB_URL` | `postgres://auth_user:auth_password@localhost:5433/auth_db?sslmode=disable` | Connection string for Auth DB |
| `JWT_SECRET` | `my_super_secret_key` | Secret key for JWT signing |

### Portfolio Service Environment Variables

*   *Currently uses default hardcoded configurations for initial setup.*

### Market Service Environment Variables

*   *Currently uses default hardcoded configurations for initial setup.*

## API Documentation

Each service exposes a `/health` endpoint to verify connectivity.

### Auth API

**Register User**
```http
POST http://localhost:8081/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Login User**
```http
POST http://localhost:8081/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```
