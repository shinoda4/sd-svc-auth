# SD-SVC-AUTH

Welcome to the documentation for **sd-svc-auth**, a robust and scalable authentication microservice written in Go.

## Overview

`sd-svc-auth` provides a complete solution for user authentication and management, designed to be easily integrated into larger distributed systems. Built with a **gRPC-first architecture**, it uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to automatically expose HTTP/JSON APIs alongside native gRPC endpoints, providing maximum flexibility for different client types.

The service handles the complexities of secure password storage, token management, and email verification, allowing you to focus on building your core business logic.

## Key Features

- **Dual Protocol Support**: Native gRPC service with automatic HTTP/JSON API via grpc-gateway
- **Secure Authentication**: Uses industry-standard JWT (JSON Web Tokens) for stateless authentication
- **User Management**: Supports user registration with email verification to ensure user validity
- **Token Rotation**: Implements access and refresh token flows to balance security and user experience
- **Password Security**: Uses strong hashing algorithms (bcrypt) for password storage
- **Password Reset**: Secure, email-based password reset flow
- **Logout Mechanism**: Token invalidation support (blacklist) for secure logout
- **Clean Architecture**: Built with maintainability and testability in mind, following standard Go project layout
- **Production Ready**: Comprehensive error handling, logging, and authentication interceptors

## Technology Stack

- **Language**: Go 1.21+
- **Protocol**: gRPC + grpc-gateway
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT (golang-jwt/jwt)
- **Configuration**: Environment variables

## Quick Start

1. **Set up dependencies**: PostgreSQL and Redis
2. **Configure environment**: Copy `.env.example` to `.env` and update values
3. **Run migrations**: Initialize the database schema
4. **Start the service**: `make run` or `go run cmd/server/main.go`

The service will start both:
- gRPC server on port `50051` (default)
- HTTP gateway on port `8080` (default)

For detailed setup instructions, see the [Deployment](./deployment.md) section.