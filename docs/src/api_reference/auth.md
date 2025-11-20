# Authentication & Tokens

This page explains the end-to-end authentication model shared by the gRPC and HTTP transports.

## Token types

| Token             | Purpose                               | Lifetime                                | Storage                             |
| ----------------- | ------------------------------------- | --------------------------------------- | ----------------------------------- |
| **Access token**  | Presented on every protected request. | `JWT_EXPIRE_HOURS` (default 1 hour).    | Not stored server-side.             |
| **Refresh token** | Used to mint a new access token.      | `JWT_REFRESH_HOURS` (default 72 hours). | Stored in Redis (`token:<userID>`). |

Both tokens are HMAC-SHA256 JWTs generated in `pkg/token`. Claims contain:

```json
{
  "token_type": "access",
  "uid": "<user-id>",
  "email": "user@example.com",
  "exp": 1732123456,
  "iat": 1732119856
}
```

## Registration & verification

1. **Register** – `POST /api/v1/register` or `AuthService.Register`. The service stores the user, creates a random verify token, and emails `SERVER_HOST:SERVER_PORT/api/v1/verify?token=<token>`.
2. **Verify** – `GET /api/v1/verify?token=...` or `AuthService.VerifyEmail`. Sets `email_verified=true`. Leave `send_email` unset/false to deliver a welcome email; set it to true to suppress the follow-up email (useful in tests).

## Login

1. Call `POST /api/v1/login` or `AuthService.Login` with email/password.
2. The service validates the bcrypt hash and ensures the email is verified.
3. On success, the refresh token is cached in Redis and both tokens are returned to the client.

### Typical HTTP response

```json
{
  "access_token": "<access>",
  "refresh_token": "<refresh>",
  "expires_in": "2025-11-20T12:00:00Z",
  "refresh_expires_in": "2025-11-23T12:00:00Z"
}
```

## Refreshing sessions

1. Send the refresh token as the Bearer credential:
   ```
   Authorization: Bearer <refresh_token>
   ```
2. Call `POST /api/v1/refresh` or `AuthService.RefreshToken`.
3. The interceptor validates the JWT, `service.Refresh` checks Redis for `token:<userID>`, and – if it matches – issues a new access token.

Refresh tokens rotate only when the user logs in again; the service simply reuses the original refresh token while it is valid.

## Logout & revocation

- Call `POST /api/v1/logout` or `AuthService.Logout` with whichever token you want to revoke.
- Access tokens are pushed to the Redis blacklist (`blacklist:<token>`). Entries inherit the original TTL, so they expire naturally.
- Refresh tokens trigger deletion of `token:<userID>` which forces users to reauthenticate before refreshing again.

## Token validation

- `POST /api/v1/verify-token` or `AuthService.ValidateToken` – Confirms that the supplied token has a valid signature, has not expired, and (for refresh tokens) matches the cached value.
- `GET /api/v1/me` or `AuthService.Me` – Returns the user ID, email, issued-at, and expiry derived from the token.

## Error scenarios

| Situation                        | Response                                                                           |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| Password mismatch                | `codes.InvalidArgument` / `400 Bad Request` with message `invalid password`.       |
| Email not verified               | `codes.FailedPrecondition` (mapped to HTTP 400) with message `email not verified`. |
| Refresh token missing from Redis | `codes.Unauthenticated` / `401 Unauthorized` with message `invalid token`.         |
| Token expired                    | `codes.Unauthenticated` / `401 Unauthorized`.                                      |
| Duplicate registration           | `codes.AlreadyExists` / `409 Conflict`.                                            |
