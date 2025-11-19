# HTTP Gateway API Reference

This page documents the HTTP/JSON API exposed via [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

## Overview

The HTTP gateway automatically translates HTTP/JSON requests to gRPC calls, providing a RESTful-style API for clients that don't support gRPC. The endpoints are defined using `google.api.http` annotations in the proto file.

## Base URL

```
http://localhost:8080
```

## Content Type

All requests and responses use JSON:

```
Content-Type: application/json
```

## Authentication

Protected endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

## Endpoints

### Register

Register a new user account.

**Endpoint**: `POST /api/v1/register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword"
}
```

**Response** (200 OK):
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "registered",
  "verify_token": "eyJhbGc..."
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "securepassword"
  }'
```

---

### Verify Email

Verify a user's email address.

**Endpoint**: `GET /api/v1/verify`

**Query Parameters**:
- `token` (required): Verification token
- `send_email` (optional): Whether to send verification email

**Response** (200 OK):
```json
{
  "message": "email verified"
}
```

**cURL Example**:
```bash
curl -X GET "http://localhost:8080/api/v1/verify?token=verification_token_here"
```

---

### Login

Authenticate and receive access and refresh tokens.

**Endpoint**: `POST /api/v1/login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": "2024-11-20T12:00:00Z",
  "refresh_expires_in": "2024-11-23T12:00:00Z"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

---

### Logout

Invalidate the current access token.

**Endpoint**: `POST /api/v1/logout`

**Authentication**: Required

**Request Body**:
```json
{}
```

**Response** (200 OK):
```json
{
  "message": "logout successful"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### Get Current User (Me)

Retrieve the authenticated user's profile.

**Endpoint**: `GET /api/v1/me`

**Authentication**: Required

**Response** (200 OK):
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "expires_in": "2024-11-20T12:00:00Z",
  "issued_at": "2024-11-19T12:00:00Z"
}
```

**cURL Example**:
```bash
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### Validate Token

Validate an access token and retrieve user information.

**Endpoint**: `POST /api/v1/verify-token`

**Authentication**: Required

**Request Body**:
```json
{}
```

**Response** (200 OK):
```json
{
  "valid": true,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/verify-token \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### Refresh Token

Obtain a new access token using a refresh token.

**Endpoint**: `POST /api/v1/refresh`

**Authentication**: Required (use refresh token)

**Request Body**:
```json
{}
```

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": "2024-11-20T12:00:00Z"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/refresh \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### Forgot Password

Request a password reset email.

**Endpoint**: `POST /api/v1/forgot-password`

**Request Body**:
```json
{
  "email": "user@example.com",
  "username": "johndoe"
}
```

**Response** (200 OK):
```json
{
  "message": "reset password email sent"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe"
  }'
```

---

### Reset Password

Complete the password reset process.

**Endpoint**: `POST /api/v1/reset-password`

**Request Body**:
```json
{
  "new_password": "newsecurepassword",
  "new_password_confirm": "newsecurepassword",
  "token": "reset_token_from_email"
}
```

**Response** (200 OK):
```json
{
  "message": "password reset done!"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "new_password": "newsecurepass",
    "new_password_confirm": "newsecurepass",
    "token": "reset_token"
  }'
```

---

## Complete Endpoint Summary

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/register` | POST | No | Register new user |
| `/api/v1/verify` | GET | No | Verify email address |
| `/api/v1/login` | POST | No | Authenticate user |
| `/api/v1/logout` | POST | Yes | Logout user |
| `/api/v1/me` | GET | Yes | Get current user |
| `/api/v1/verify-token` | POST | Yes | Validate token |
| `/api/v1/refresh` | POST | Yes (refresh) | Refresh access token |
| `/api/v1/forgot-password` | POST | No | Request password reset |
| `/api/v1/reset-password` | POST | No | Complete password reset |

## Error Responses

Errors are returned with appropriate HTTP status codes and JSON body:

**Example Error** (401 Unauthorized):
```json
{
  "code": 16,
  "message": "invalid token",
  "details": []
}
```

### HTTP Status Code Mapping

| gRPC Code | HTTP Status | Description |
|-----------|-------------|-------------|
| `OK` | 200 | Success |
| `INVALID_ARGUMENT` | 400 | Bad Request |
| `UNAUTHENTICATED` | 401 | Unauthorized |
| `PERMISSION_DENIED` | 403 | Forbidden |
| `NOT_FOUND` | 404 | Not Found |
| `ALREADY_EXISTS` | 409 | Conflict |
| `INTERNAL` | 500 | Internal Server Error |

## JavaScript/TypeScript Example

```typescript
// Register
const registerResponse = await fetch('http://localhost:8080/api/v1/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'user@example.com',
    username: 'johndoe',
    password: 'securepassword',
  }),
});

const registerData = await registerResponse.json();
console.log('User ID:', registerData.user_id);

// Login
const loginResponse = await fetch('http://localhost:8080/api/v1/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'securepassword',
  }),
});

const { access_token, refresh_token } = await loginResponse.json();

// Get current user
const meResponse = await fetch('http://localhost:8080/api/v1/me', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${access_token}`,
  },
});

const user = await meResponse.json();
console.log('User:', user);

// Refresh token
const refreshResponse = await fetch('http://localhost:8080/api/v1/refresh', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${refresh_token}`,
  },
  body: '{}',
});

const { access_token: newAccessToken } = await refreshResponse.json();
```

## Python Example

```python
import requests

# Register
register_response = requests.post(
    'http://localhost:8080/api/v1/register',
    json={
        'email': 'user@example.com',
        'username': 'johndoe',
        'password': 'securepassword'
    }
)
user_data = register_response.json()
print(f"User ID: {user_data['user_id']}")

# Login
login_response = requests.post(
    'http://localhost:8080/api/v1/login',
    json={
        'email': 'user@example.com',
        'password': 'securepassword'
    }
)

tokens = login_response.json()
access_token = tokens['access_token']
refresh_token = tokens['refresh_token']

# Get current user
me_response = requests.get(
    'http://localhost:8080/api/v1/me',
    headers={'Authorization': f'Bearer {access_token}'}
)

user = me_response.json()
print(f"User: {user}")

# Refresh token
refresh_response = requests.post(
    'http://localhost:8080/api/v1/refresh',
    headers={'Authorization': f'Bearer {refresh_token}'},
    json={}
)

new_tokens = refresh_response.json()
new_access_token = new_tokens['access_token']
```

## Notes

- All endpoints use the `/api/v1/` prefix as defined in the proto file
- `GET` methods (`/verify` and `/me`) don't require a request body
- Query parameters for `GET /api/v1/verify` are automatically mapped from the proto message fields
- The HTTP gateway respects the `google.api.http` annotations in the proto file

