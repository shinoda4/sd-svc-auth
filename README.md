# SD-SVC-AUTH

A robust and scalable authentication microservice written in Go, providing gRPC and HTTP/JSON APIs for user authentication and management.

## Overview

`sd-svc-auth` is a production-ready authentication service built with a **gRPC-first architecture**. It uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to automatically expose HTTP/JSON APIs alongside native gRPC endpoints, providing flexibility for different client types.

## Key Features

- **Dual Protocol Support**: Native gRPC + HTTP/JSON APIs via grpc-gateway
- **Secure Authentication**: JWT-based stateless authentication with access and refresh tokens
- **User Management**: Registration with email verification
- **Token Management**: Token validation, refresh, and blacklist-based logout
- **Password Security**: bcrypt hashing with secure password reset flow
- **Clean Architecture**: Separation of concerns with repository, service, and transport layers
- **Production Ready**: Comprehensive error handling, logging, and authentication interceptors

## Architecture

The project follows the **Standard Go Project Layout** with clean architecture principles:

```
sd-svc-auth/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── model/            # Domain models
│   ├── repo/             # Data access layer (PostgreSQL + Redis)
│   ├── service/          # Business logic layer
│   │   └── auth/         # Authentication service
│   └── transport/        # Transport layer
│       └── grpc/         # gRPC server + HTTP gateway
├── pkg/
│   ├── email/            # Email sending utility
│   ├── logger/           # Logging utility
│   └── token/            # JWT generation and validation
├── db/                   # Database migrations
├── docs/                 # Documentation (mdbook)
└── deployments/          # Docker and deployment configs
```

### Transport Layer

- **gRPC Server**: Native gRPC service implementing `auth.v1.AuthService`
- **HTTP Gateway**: Automatic HTTP/JSON API via grpc-gateway
- **Authentication Interceptor**: Validates JWT tokens for protected endpoints

## Technology Stack

- **Language**: Go 1.21+
- **Protocol**: gRPC + grpc-gateway
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT (golang-jwt/jwt)
- **Email**: SMTP

## Setup and Installation

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Redis 6 or higher

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Database
DATABASE_DSN=postgres://sd_auth:sd_pass@localhost:5432/sd_auth?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# Server
SERVER_HOST=localhost
SERVER_PORT=8080
HTTP_PORT=8080      # HTTP gateway port
GRPC_PORT=50051     # gRPC server port

# JWT
JWT_SECRET=replace_with_a_strong_secret
JWT_EXPIRE_HOURS=72

# Email (SMTP)
EMAIL_ADDRESS=your_email@example.com
EMAIL_PASSWORD=your_app_password

# Password Reset
RESET_PASSWORD_URL=http://localhost:8080/api/v1/reset-password
```

### Running the Service

#### Using Make

```bash
# Build the binary
make build

# Run the service
make run

# Run tests
make test

# Start with Docker Compose
make up
```

#### Direct Execution

```bash
# Install dependencies
go mod download

# Run the service
go run cmd/server/main.go
```

The service will start:

- **gRPC server** on port `50051` (default)
- **HTTP gateway** on port `8080` (default)

## API Documentation

The service provides both gRPC and HTTP/JSON APIs.

### gRPC API

**Service**: `auth.v1.AuthService`

**Methods**:

- `Register` - Register a new user
- `Login` - Authenticate and receive tokens
- `VerifyEmail` - Verify email address
- `Logout` - Invalidate access token
- `RefreshToken` - Refresh access token
- `ValidateToken` - Validate token
- `Me` - Get current user profile
- `ForgotPassword` - Request password reset
- `ResetPassword` - Complete password reset

### HTTP/JSON API

**Base URL**: `http://localhost:8080`

#### Authentication Endpoints

- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - Login and receive tokens
- `GET /api/v1/verify` - Verify email address
- `POST /api/v1/logout` - Logout (requires auth)
- `POST /api/v1/refresh` - Refresh access token (requires auth)
- `POST /api/v1/verify-token` - Validate token (requires auth)

#### User Endpoints

- `GET /api/v1/me` - Get current user (requires auth)

#### Password Reset Endpoints

- `POST /api/v1/forgot-password` - Request password reset
- `POST /api/v1/reset-password` - Complete password reset

### Authentication

Protected endpoints require a Bearer token in the `Authorization` header or gRPC metadata:

**HTTP**:

```
Authorization: Bearer <access_token>
```

**gRPC**:

```
authorization: Bearer <access_token>
```

## Database Setup

### Run Migrations

```bash
# Using golang-migrate
migrate -source file://db/migrations -database "$DATABASE_DSN" up

# Or using Make (if configured)
make init-db
```

## Development

### Project Structure

- **Transport Layer** (`internal/transport/grpc/`): gRPC handlers and HTTP gateway
- **Service Layer** (`internal/service/auth/`): Business logic
- **Repository Layer** (`internal/repo/`): Database and cache operations
- **Models** (`internal/model/`): Domain entities

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./tests/...
```

## Documentation

Full documentation is available in the `docs/` directory using [mdbook](https://rust-lang.github.io/mdBook/).

### Build Documentation

```bash
cd docs
mdbook build
mdbook serve  # Serve at http://localhost:3000
```

## Docker Deployment

```bash
# Build image
docker build -t sd-svc-auth:latest .

# Run container
docker run -p 8080:8080 -p 50051:50051 --env-file .env sd-svc-auth:latest

# Or use Docker Compose
docker-compose -f deployments/docker-compose.yml up
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
