# Architecture

`sd-svc-auth` follows the Standard Go Project Layout and applies Clean Architecture principles: external transports stay thin, business rules live in the service layer, and infrastructure (database, cache, email, token utilities) is pushed behind interfaces.

## Layered overview

```
cmd/server/main.go
        |
        v
internal/config.MustLoad()   pkg/logger.Init()
        |
        v
internal/repo    <-- PostgreSQL + Redis adapters
        |
        v
internal/service/auth        <-- business logic
        |
        v
internal/transport/grpc      <-- gRPC server + HTTP gateway
```

- **cmd/server** wires configuration, repositories, and services, starts the gRPC server and HTTP gateway, and handles graceful shutdown.
- **internal/config** validates required environment variables before the process advertises any listener.
- **internal/repo** provides concrete implementations of `entity.UserRepository` (PostgreSQL via `sqlx`) and `entity.CacheRepository` (Redis). Password hashing uses bcrypt.
- **internal/service/auth** hosts every use-case (register, verify email, login, token refresh, logout, password reset, token validation). The layer is written against the repository interfaces.
- **internal/transport/grpc** exposes the service over gRPC, adds interceptors (logging + authentication), registers the grpc-gateway HTTP handler, and implements a lightweight health probe.
- **pkg** holds shared helper packages (`token`, `email`, `logger`) that do not depend on application internals.

## Repository layer

- **PostgreSQL (`internal/repo/user_repo.go`)**
  - Prevents duplicate registrations by checking `users.email`.
  - Persists password hashes, verification tokens, and password-reset metadata.
  - Provides helpers such as `GetUserByVerifyToken`, `SaveResetToken`, and `ClearResetToken`.
- **Redis (`internal/repo/redis_cache.go`)**
  - Stores refresh tokens (`token:<userID>`), enforces refresh-token reuse detection, and deletes entries on logout.
  - Maintains an access-token blacklist (`blacklist:<token>`) for immediate revocation.

## Service layer

Each method in `internal/service/auth` accepts a `context.Context` and coordinates repositories + helpers:

- **Register** – Creates the user, generates a verification token, stores it, and optionally sends an email through `pkg/email`.
- **VerifyEmail** – Validates the token, marks the user verified, and can send a welcome email.
- **Login** – Validates credentials, enforces email verification, issues an access/refresh token pair via `pkg/token`, and caches the refresh token in Redis.
- **Refresh** – Validates the refresh token, ensures it matches the cached value, and issues a new access token.
- **Logout** – Adds access tokens to the blacklist or clears refresh tokens, depending on the token type.
- **PasswordReset**/**PasswordResetConfirm** – Issues random reset tokens, emails reset links using `RESET_PASSWORD_URL`, updates the password, then clears reset secrets.
- **ValidateToken**/**Me** – Helper endpoints that surface the JWT claims for clients.

## Transport layer

### gRPC server

- Implements `authpb.AuthServiceServer` (generated from `github.com/shinoda4/sd-grpc-proto/proto/auth/v1`).
- Listens on `:$GRPC_PORT` and shares a chain of interceptors:
  - **Logging interceptor** – emits the method name and error (if any).
  - **Auth interceptor** – skips public RPCs (`Register`, `Login`, `VerifyEmail`, `ForgotPassword`, `ResetPassword`, `HealthCheck`) and enforces Bearer tokens everywhere else. Valid JWT claims are injected into the context under `claims`.
- Provides a lightweight `HealthCheck` RPC that returns `"ok"` and is whitelisted from authentication.

### HTTP gateway

- Built with `grpc-gateway/runtime.ServeMux`.
- Dials the gRPC server locally and forwards requests using the `google.api.http` annotations defined in the proto.
- Listens on `:$HTTP_PORT` (default 8080) and exposes REST endpoints under `/api/v1/*`.

## Data flow examples

### Login

1. Client issues `Login` (gRPC or HTTP) with email/password.
2. Transport maps the request to `service.Login`.
3. Service fetches the user, checks bcrypt hash, verifies `email_verified`, generates tokens, caches refresh token.
4. Transport converts the response into Protobuf timestamps / JSON datetimes.

### Refresh

1. Client calls `RefreshToken` with a Bearer refresh token.
2. Auth interceptor validates the JWT and attaches claims.
3. `service.Refresh` ensures the supplied token matches the cached version, then issues a new access token.

### Password reset

1. `ForgotPassword` validates the username/email pair and writes a random 64-char hex token plus an expiry to PostgreSQL.
2. An email is sent via SMTP using `EMAIL_ADDRESS`/`EMAIL_PASSWORD`.
3. `ResetPassword` checks the token, ensures it has not expired, updates the stored hash, and clears reset metadata.

## Error handling

- The service layer returns typed errors (`service.ErrInvalidPassword`, `service.ErrEmailNotVerified`, etc.). The transport layer maps them to gRPC status codes, which grpc-gateway translates to HTTP responses.
- Repository errors are wrapped with context (e.g., `insert user`, `query user`) to simplify troubleshooting.
- JWT parsing/validation errors surface as `codes.Unauthenticated`.

## Observability hooks

- `pkg/logger.Init()` configures standard library logging with timestamps + file names.
- The gRPC logging interceptor reports method names and errors.
- Docker/Kubernetes deployments can hook into the `HealthCheck` RPC or add an HTTP `/health` handler in front of the gateway if needed.
