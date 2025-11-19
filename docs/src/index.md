# SD-SVC-AUTH

Welcome to the documentation for **sd-svc-auth**, a robust and scalable authentication microservice written in Go.

## Overview

`sd-svc-auth` provides a complete solution for user authentication and management, designed to be easily integrated into larger distributed systems. It handles the complexities of secure password storage, token management, and email verification, allowing you to focus on building your core business logic.

## Key Features

- **Secure Authentication**: Uses industry-standard JWT (JSON Web Tokens) for stateless authentication.
- **User Management**: Supports user registration with email verification to ensure user validity.
- **Token Rotation**: Implements access and refresh token flows to balance security and user experience.
- **Password Security**: Uses strong hashing algorithms (bcrypt) for password storage.
- **Password Reset**: Secure, email-based password reset flow.
- **Logout Mechanism**: Token invalidation support (blacklist) for secure logout.
- **Clean Architecture**: Built with maintainability and testability in mind, following standard Go project layout.

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT (golang-jwt/jwt)
- **Configuration**: Viper (or standard env vars)