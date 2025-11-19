# gRPC API Reference

This page documents the native gRPC API for the authentication service.

## Service Definition

**Package**: `auth.v1`  
**Service**: `AuthService`  
**Proto Repository**: [sd-grpc-proto](https://github.com/shinoda4/sd-grpc-proto)

## Connection Details

- **Protocol**: gRPC (HTTP/2)
- **Default Port**: `50051`
- **TLS**: Not enabled by default (use insecure credentials for local development)

## Authentication

Most endpoints require authentication via gRPC metadata:

```
authorization: Bearer <access_token>
```

### Public Endpoints (No Authentication Required)

- `Register`
- `Login`
- `VerifyEmail`
- `ForgotPassword`
- `ResetPassword`

### Protected Endpoints (Authentication Required)

- `Logout`
- `RefreshToken`
- `ValidateToken`
- `Me`

## Methods

### Register

Register a new user account with email verification.

**RPC**: `Register(RegisterRequest) returns (RegisterResponse)`

**Request**:
```protobuf
message RegisterRequest {
  string email = 1;
  string username = 2;
  string password = 3;
  bool send_email = 4;  // Optional, defaults to true
}
```

**Response**:
```protobuf
message RegisterResponse {
  string user_id = 1;
  string message = 2;
  string verify_token = 3;  // Only for testing, normally sent via email
}
```

**Example** (Go client):
```go
resp, err := client.Register(ctx, &authpb.RegisterRequest{
    Email:    "user@example.com",
    Username: "johndoe",
    Password: "securepassword",
})
```

---

### Login

Authenticate a user and receive access and refresh tokens.

**RPC**: `Login(LoginRequest) returns (LoginResponse)`

**Request**:
```protobuf
message LoginRequest {
  string email = 1;
  string password = 2;
}
```

**Response**:
```protobuf
message LoginResponse {
  string access_token = 1;
  string refresh_token = 2;
  google.protobuf.Timestamp expires_in = 3;
  google.protobuf.Timestamp refresh_expires_in = 4;
}
```

**Example**:
```go
resp, err := client.Login(ctx, &authpb.LoginRequest{
    Email:    "user@example.com",
    Password: "securepassword",
})
```

---

### VerifyEmail

Verify a user's email address using the verification token.

**RPC**: `VerifyEmail(VerifyEmailRequest) returns (VerifyEmailResponse)`

**Request**:
```protobuf
message VerifyEmailRequest {
  string token = 1;
  bool send_email = 2;  // If true, sends verification email
}
```

**Response**:
```protobuf
message VerifyEmailResponse {
  string message = 1;
}
```

---

### Logout

Invalidate the current access token (adds to blacklist).

**RPC**: `Logout(LogoutRequest) returns (LogoutResponse)`

**Authentication**: Required

**Request**:
```protobuf
message LogoutRequest {}
```

**Response**:
```protobuf
message LogoutResponse {
  string message = 1;
}
```

**Example**:
```go
md := metadata.Pairs("authorization", "Bearer "+accessToken)
ctx := metadata.NewOutgoingContext(context.Background(), md)

resp, err := client.Logout(ctx, &authpb.LogoutRequest{})
```

---

### RefreshToken

Obtain a new access token using a valid refresh token.

**RPC**: `RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse)`

**Authentication**: Required (refresh token in metadata)

**Request**:
```protobuf
message RefreshTokenRequest {}
```

**Response**:
```protobuf
message RefreshTokenResponse {
  string access_token = 1;
  google.protobuf.Timestamp expires_in = 2;
}
```

---

### ValidateToken

Validate an access token and retrieve user information.

**RPC**: `ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)`

**Authentication**: Required

**Request**:
```protobuf
message ValidateTokenRequest {}
```

**Response**:
```protobuf
message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  string email = 3;
}
```

---

### Me

Get the current authenticated user's profile information.

**RPC**: `Me(MeRequest) returns (MeResponse)`

**Authentication**: Required

**Request**:
```protobuf
message MeRequest {}
```

**Response**:
```protobuf
message MeResponse {
  string user_id = 1;
  string email = 2;
  google.protobuf.Timestamp expires_in = 3;
  google.protobuf.Timestamp issued_at = 4;
}
```

---

### ForgotPassword

Initiate the password reset flow by sending a reset email.

**RPC**: `ForgotPassword(ForgotPasswordRequest) returns (ForgotPasswordResponse)`

**Request**:
```protobuf
message ForgotPasswordRequest {
  string email = 1;
  string username = 2;
}
```

**Response**:
```protobuf
message ForgotPasswordResponse {
  string message = 1;
}
```

---

### ResetPassword

Complete the password reset using the reset token.

**RPC**: `ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse)`

**Request**:
```protobuf
message ResetPasswordRequest {
  string token = 1;
  string new_password = 2;
  string new_password_confirm = 3;
}
```

**Response**:
```protobuf
message ResetPasswordResponse {
  string message = 1;
}
```

## Error Handling

The service uses standard gRPC status codes:

| Code | Description | Example Scenario |
|------|-------------|------------------|
| `OK` | Success | Successful operation |
| `INVALID_ARGUMENT` | Invalid input | Password mismatch, missing fields |
| `UNAUTHENTICATED` | Authentication failed | Invalid token, missing credentials |
| `ALREADY_EXISTS` | Resource exists | Email already registered |
| `NOT_FOUND` | Resource not found | User doesn't exist |
| `INTERNAL` | Server error | Database connection failure |

**Example Error Handling** (Go):
```go
resp, err := client.Login(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.Unauthenticated:
            // Handle invalid credentials
        case codes.InvalidArgument:
            // Handle validation errors
        }
    }
}
```

## Client Examples

### Go Client

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    authpb "github.com/shinoda4/sd-grpc-proto/auth/v1"
)

// Connect to server
conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := authpb.NewAuthServiceClient(conn)

// Register
resp, err := client.Register(context.Background(), &authpb.RegisterRequest{
    Email:    "user@example.com",
    Username: "johndoe",
    Password: "securepass",
})
```

### Python Client

```python
import grpc
from auth.v1 import auth_pb2, auth_pb2_grpc

# Connect to server
channel = grpc.insecure_channel('localhost:50051')
client = auth_pb2_grpc.AuthServiceStub(channel)

# Login
response = client.Login(auth_pb2.LoginRequest(
    email='user@example.com',
    password='securepass'
))

print(f"Access Token: {response.access_token}")
```
