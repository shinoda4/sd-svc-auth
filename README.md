# Auth Service

## Overview
This is a robust authentication microservice written in Go, designed to handle user registration, login, token management (JWT), and password resets. It follows a clean architecture pattern to ensure maintainability and scalability.

## Architecture
The project is structured following the Standard Go Project Layout:

- **`cmd/`**: Main applications for this project.
- **`internal/`**: Private application and library code.
    - **`handler/`**: HTTP handlers for processing requests.
    - **`service/`**: Business logic layer.
    - **`repo/`**: Data access layer (database and cache).
    - **`dto/`**: Data Transfer Objects for API requests and responses.
    - **`model/`**: Domain models.
- **`pkg/`**: Library code that's ok to use by external applications (e.g., token utilities, email sender).
- **`tests/`**: Integration and unit tests.

## Features
- **User Registration**: With email verification.
- **Login**: JWT-based authentication (Access & Refresh Tokens).
- **Token Management**: Refresh token rotation and validation.
- **Password Reset**: Secure password reset flow via email.
- **Logout**: Token invalidation (blacklist).

## Setup and Run

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis

### Environment Variables
Create a `.env` file or set the following environment variables:
```bash
DB_SOURCE=postgresql://user:password@localhost:5432/auth_db?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_SECRET=your_secret_key
EMAIL_ADDRESS=your_email@example.com
EMAIL_PASSWORD=your_email_password
```

### Running the Service
```bash
go run cmd/server/main.go
```

### Running Tests
```bash
go test ./tests/...
```

## API Documentation

### Auth
- `POST /api/v1/register`: Register a new user.
- `POST /api/v1/login`: Login and receive tokens.
- `POST /api/v1/refresh`: Refresh access token.
- `POST /api/v1/logout`: Logout user.
- `POST /api/v1/verify-token`: Verify validity of an access token.
- `GET /api/v1/verify`: Verify email address.

### Password Reset
- `POST /api/v1/password-reset`: Request password reset.
- `POST /api/v1/password-reset-confirm`: Confirm new password.

### User
- `GET /api/v1/authorized/me`: Get current user profile.