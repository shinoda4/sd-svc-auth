# gRPC API Reference

All RPCs live in the `auth.v1` package and are defined in [`github.com/shinoda4/sd-grpc-proto/proto/auth/v1`](https://github.com/shinoda4/sd-grpc-proto). The server listens on `:$GRPC_PORT` (default `50051`) and currently exposes an insecure endpoint for local development—add TLS before running in production.

```go
conn, err := grpc.Dial(
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
client := authpb.NewAuthServiceClient(conn)
```

Authentication is provided by a standard Bearer token embedded in gRPC metadata:

```go
md := metadata.Pairs("authorization", "Bearer "+accessToken)
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

Public methods that do not require authentication: `HealthCheck`, `Register`, `Login`, `VerifyEmail`, `ForgotPassword`, `ResetPassword`.

## Methods

### HealthCheck

```protobuf
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
```

Returns a static `"ok"` status. Use it for readiness probes.

### Register

```protobuf
message RegisterRequest {
  string email = 1;
  string username = 2;
  string password = 3;
  bool send_email = 4; // optional, default true (verification email is sent)
}
```

```protobuf
message RegisterResponse {
  string user_id = 1;
  string message = 2;
  string verify_token = 3; // only surfaced for testing
}
```

Generates a verification token, stores it in PostgreSQL, and sends an email using `SERVER_HOST`/`SERVER_PORT` to build the link.

### Login

```protobuf
rpc Login(LoginRequest) returns (LoginResponse);
```

Returns both access and refresh tokens. Email addresses must already be verified.

```protobuf
message LoginResponse {
  string access_token = 1;
  string refresh_token = 2;
  google.protobuf.Timestamp expires_in = 3;
  google.protobuf.Timestamp refresh_expires_in = 4;
}
```

### VerifyEmail

```protobuf
message VerifyEmailRequest {
  string token = 1;
  bool send_email = 2; // set to false (default) to send a welcome email
}
```

The service validates the `verify_token` stored in the database. Passing `send_email=true` suppresses the follow-up welcome email (useful for integration tests).

### ForgotPassword

```protobuf
message ForgotPasswordRequest {
  string email = 1;
  string username = 2;
}
```

Generates a hex-encoded reset token, stores it with a one-hour expiry, and emails `RESET_PASSWORD_URL?token=<...>` to the user.

### ResetPassword

```protobuf
message ResetPasswordRequest {
  string token = 1;
  string new_password = 2;
  string new_password_confirm = 3;
}
```

Ensures the token exists, has not expired, and that both passwords match. After success, the token is cleared.

### Logout

```protobuf
rpc Logout(LogoutRequest) returns (LogoutResponse); // auth required
```

If the supplied token is an access token, it is placed on a Redis blacklist for the remainder of its TTL. If it is a refresh token, the cached refresh token is deleted.

### RefreshToken

```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse); // requires Bearer refresh token
```

Validates the refresh token against the cached value in Redis and issues a new access token.

### ValidateToken

```protobuf
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse); // auth required
```

Returns the parsed token claims (`user_id`, `email`, and `valid=true`) if the supplied token is still active and not blacklisted.

### Me

```protobuf
rpc Me(MeRequest) returns (MeResponse); // auth required
```

Reads the claims injected by the interceptor and returns:

```protobuf
message MeResponse {
  string user_id = 1;
  string email = 2;
  google.protobuf.Timestamp expires_in = 3;
  google.protobuf.Timestamp issued_at = 4;
}
```

## Error handling

Use `status.FromError(err)` to inspect gRPC codes:

```go
resp, err := client.Login(ctx, req)
if err != nil {
    if st, ok := status.FromError(err); ok {
        switch st.Code() {
        case codes.Unauthenticated:
            // invalid credentials or token
        case codes.InvalidArgument:
            // validation issues
        }
    }
}
```

## Metadata summary

| Method | Auth required | Notes |
|--------|---------------|-------|
| `HealthCheck` | No | Probing |
| `Register`, `Login`, `VerifyEmail`, `ForgotPassword`, `ResetPassword` | No | Public entry points |
| `Logout`, `RefreshToken`, `ValidateToken`, `Me` | Yes | Requires Bearer token |
