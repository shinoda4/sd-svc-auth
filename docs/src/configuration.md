# Configuration

`sd-svc-auth` is configured entirely through environment variables. The process fails fast during startup if any required variable is missing.

## Core environment variables

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DATABASE_DSN` | ✅ | PostgreSQL DSN used by `sqlx`. | `postgres://user:pass@localhost:5432/sd_auth?sslmode=disable` |
| `REDIS_ADDR` | ✅ | Redis host:port for refresh tokens + blacklist. | `localhost:6379` |
| `REDIS_PASSWORD` | ❌ | Redis password, leave empty if not set. | `supersecret` |
| `SERVER_HOST` | ✅ | Base host used to build verification URLs. | `http://localhost` or `https://api.example.com` |
| `SERVER_PORT` | ✅ | Port appended to `SERVER_HOST` for verification links. | `8080` |
| `HTTP_PORT` | ✅ | Port for the grpc-gateway HTTP server. | `8080` |
| `GRPC_PORT` | ✅ | Port for the gRPC server. | `50051` |
| `JWT_SECRET` | ✅ | HMAC secret for signing JWTs. Must be ≥32 bytes. | `openssl rand -base64 32` |
| `JWT_EXPIRE_HOURS` | ❌ | Access token lifetime in hours (default 1). | `72` |
| `JWT_REFRESH_HOURS` | ❌ | Refresh token lifetime in hours (default 72). | `168` |
| `EMAIL_ADDRESS` | ✅ | SMTP username / from-address. | `noreply@example.com` |
| `EMAIL_PASSWORD` | ✅ | SMTP password or app password. | `app-specific-pass` |
| `RESET_PASSWORD_URL` | ✅ | Base URL used in reset emails (`?token=` is appended). | `https://app.example.com/reset-password` |

> **Note:** `SERVER_HOST` + `SERVER_PORT` are only used to build email verification links. HTTP traffic for users still goes through `$HTTP_PORT`.

## Loading configuration

`internal/config.MustLoad()` reads the variables above, verifies required entries, and returns a struct that is injected into repositories and transports. Missing variables trigger `log.Fatalf`, preventing partially configured nodes from accepting traffic.

## Example `.env` for local development

```bash
# Database
DATABASE_DSN=postgres://sd_auth:sd_pass@localhost:5432/sd_auth?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# Network
SERVER_HOST=http://localhost
SERVER_PORT=8080
HTTP_PORT=8080
GRPC_PORT=50051

# JWT
JWT_SECRET=dev-secret-change-this
JWT_EXPIRE_HOURS=24
JWT_REFRESH_HOURS=168

# Email
EMAIL_ADDRESS=your-email@gmail.com
EMAIL_PASSWORD=your-app-password

# Reset links
RESET_PASSWORD_URL=http://localhost:8080/api/v1/reset-password
```

Load the file during development with `export $(cat .env | xargs)` or by using a process manager that reads `.env`.

## Production hints

- Use managed secret stores (AWS Secrets Manager, HashiCorp Vault, Kubernetes secrets) instead of bundling credentials in images.
- Set `sslmode=require` (or stronger) in `DATABASE_DSN` and put PostgreSQL behind TLS.
- Store `JWT_SECRET` in HSM-backed services or rotate it periodically.
- Use a dedicated SMTP account for `EMAIL_ADDRESS` and lock it down to send-only credentials.
- When running behind TLS-terminating load balancers, set `SERVER_HOST` to the full HTTPS origin so verification links point to the correct domain.

## Docker & Compose

`deployments/docker-compose.yml` shows how to pass environment variables into the container:

```yaml
services:
  sd-svc-auth:
    image: shinoda4/sd-svc-auth:latest
    environment:
      DATABASE_DSN: postgres://sd_svc_auth_user:sd_svc_auth_password@postgres:5432/sd_svc_auth?sslmode=disable
      REDIS_ADDR: redis:6379
      SERVER_HOST: http://0.0.0.0
      HTTP_PORT: 8080
      GRPC_PORT: 50051
      JWT_SECRET: replace_with_a_strong_secret
      EMAIL_ADDRESS: noreply@example.com
      EMAIL_PASSWORD: smtp_app_password
      RESET_PASSWORD_URL: https://app.example.com/reset-password
```

## Troubleshooting

- **Missing variables:** Startup will exit with `missing required environment variables: [...]`. Export them (or add them to your container manifest) and restart.
- **Wrong DSN:** `psql "$DATABASE_DSN"` is the fastest way to verify credentials and network access.
- **Redis auth failures:** Ensure `REDIS_PASSWORD` is set both on the server and in the environment. A blank string attempts unauthenticated access.
