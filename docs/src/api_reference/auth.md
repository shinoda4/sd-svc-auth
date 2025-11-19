# Authentication

## Register

Register a new user.

- **URL**: `/register`
- **Method**: `POST`
- **Query Params**:
  - `sendEmail` (optional, default `true`): Whether to send a verification email.

### Request Body

```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword"
}
```

### Response

```json
{
  "message": "registered",
  "verifyToken": "..." // Only returned for testing/debugging if needed, usually sent via email
}
```

## Login

Login with email and password to receive access and refresh tokens.

- **URL**: `/login`
- **Method**: `POST`

### Request Body

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

### Response

```json
{
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "expires_in": 3600,
  "refresh_expires_in": 259200
}
```

## Refresh Token

Get a new access token using a valid refresh token.

- **URL**: `/refresh`
- **Method**: `POST`

### Request Body

```json
{
  "refresh_token": "eyJhbG..."
}
```

### Response

```json
{
  "access_token": "eyJhbG...",
  "expires_in": 3600
}
```

## Verify Token

Check if an access token is valid.

- **URL**: `/verify-token`
- **Method**: `POST`

### Request Body

```json
{
  "token": "eyJhbG..."
}
```

### Response

```json
{
  "token": {
    "token_type": "access",
    "uid": "...",
    "email": "...",
    "exp": ...
  }
}
```

## Logout

Invalidate the current access token.

- **URL**: `/logout`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`

### Response

```json
{
  "message": "logout successful"
}
```
