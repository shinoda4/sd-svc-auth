# Password Reset Flow

Password resets involve two steps: requesting a reset email and confirming the reset with the emailed token. Both steps are exposed over gRPC and HTTP.

## 1. Request a reset token

### HTTP

- **Endpoint**: `POST /api/v1/forgot-password`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "username": "johndoe"
  }
  ```

### gRPC

```protobuf
rpc ForgotPassword(ForgotPasswordRequest) returns (ForgotPasswordResponse);
```

```protobuf
message ForgotPasswordRequest {
  string email = 1;
  string username = 2;
}
```

The service validates the username/email pair, generates a 64-character hex token, stores it with a one-hour expiry, and emails `RESET_PASSWORD_URL?token=<token>` via SMTP.

## 2. Confirm the reset

### HTTP

- **Endpoint**: `POST /api/v1/reset-password`
- **Body**:
  ```json
  {
    "token": "0d90f5...",
    "new_password": "MyN3wPass!",
    "new_password_confirm": "MyN3wPass!"
  }
  ```

### gRPC

```protobuf
rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse);
```

The service ensures:

1. The token exists in PostgreSQL.
2. `reset_token_expire` is still in the future.
3. `new_password` matches `new_password_confirm`.

If all checks pass, the password hash is updated (bcrypt) and both `reset_token` and `reset_token_expire` are cleared.

## Handling errors

| Error | Meaning |
|-------|---------|
| `codes.InvalidArgument` / `400 Bad Request` | Missing token, mismatched passwords, or malformed payload. |
| `codes.NotFound` / `404 Not Found` | Unknown or already used token. |
| `codes.Unauthenticated` / `401 Unauthorized` | Token expired (reported as `Token expired!`). |
| `codes.Internal` / `500 Internal Server Error` | Database/email failures. |

## Tips

- `RESET_PASSWORD_URL` should point to a frontend page (or Postman collection) that can finish the flow by calling `/api/v1/reset-password`.
- If you send multiple requests in quick succession, only the latest token remains valid because each call overwrites the previous token in PostgreSQL.
- SMTP credentials are loaded from `EMAIL_ADDRESS` / `EMAIL_PASSWORD`; ensure the account allows sending via `smtp.gmail.com:587` if you use Gmail.
