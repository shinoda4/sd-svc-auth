# SD-SVC-AUTH

`sd-svc-auth` is a standalone authentication microservice written in Go 1.25.4. It exposes a gRPC-first API and mirrors every method through an auto-generated HTTP/JSON gateway, so it can serve both internal services and public clients.

## Why use it

- **Unified transports** – Native gRPC (`auth.v1.AuthService`) with a grpc-gateway powered REST facade on `/api/v1/*`.
- **Complete auth flows** – Registration with email verification, password-based login, refreshable JWTs, token validation, logout/blacklist, and password reset.
- **Battle-tested storage** – PostgreSQL (primary store) plus Redis for refresh-token storage and token blacklist.
- **Operational extras** – Health probes, structured logging hooks, graceful shutdown, and mdBook documentation.
- **Clean layering** – `cmd` (bootstrap) → `internal/config|repo|service|transport` → `pkg` helpers (token, email, logger).

## Component overview

```
sd-svc-auth/
├── cmd/server          # Entry point that wires config, repos, services, transports
├── internal
│   ├── config          # Environment-backed configuration loader
│   ├── repo            # PostgreSQL + Redis repositories
│   ├── service/auth    # Business logic (register/login/refresh/reset/etc.)
│   └── transport/grpc  # gRPC server, gateway, auth interceptor, health probe
├── pkg
│   ├── token           # JWT + verification/reset token helpers
│   ├── email           # SMTP helper (gomail)
│   └── logger          # Logger bootstrap
├── db/migrations       # SQL migrations managed by golang-migrate
├── deployments         # Docker Compose + runtime scripts
└── docs                # This mdBook
```

## Quick start checklist

1. **Install dependencies** – Go 1.25.4+, PostgreSQL 17+, Redis 7+, `golang-migrate`, `mdbook` (optional for docs).
2. **Configure environment** – Copy `.env.example` to `.env`, fill in database/redis/JWT/email settings, and export them (see [Configuration](./configuration.md)).
3. **Run migrations** – `migrate -source file://db/migrations -database "$DATABASE_DSN" up` or `make init-db`.
4. **Start the service**
   ```bash
   make run             # local dev
   make docker-build    # build container image
   docker-compose -f deployments/docker-compose.yml up
   ```
5. **Call the APIs** – gRPC on `localhost:50051`, HTTP gateway on `http://localhost:8080/api/v1`.

## Documentation map

- [Architecture](./architecture.md) – Deep dive into layers, data flow, and interceptors.
- [Configuration](./configuration.md) – Environment variables, defaults, and production tips.
- [Deployment](./deployment.md) – Local dev, Makefile targets, Docker/Kubernetes guidance.
- [Database](./database.md) – Schema, migrations, and operational considerations.
- [API Reference](./api_reference/index.md) – gRPC/HTTP specs, examples, and auth flows.
