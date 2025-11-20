# Deployment

This guide walks through running `sd-svc-auth` locally, in containers, and in Kubernetes. All commands assume you are in the project root (`/Users/carl/sd-system/sd-svc-auth`).

## Prerequisites

- Go **1.25.4** or newer
- PostgreSQL **17** (or compatible)
- Redis **7**
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI
- [`mdbook`](https://rust-lang.github.io/mdBook/) (optional, for docs)
- Docker + Docker Compose (optional)

## Environment

1. Copy the sample file:
   ```bash
   cp .env.example .env
   ```
2. Fill in database, Redis, JWT, and email secrets (see [Configuration](./configuration.md)).
3. Export the variables before running the service:
   ```bash
   export $(grep -v '^#' .env | xargs)
   ```

## Database migrations

Run migrations anytime the schema changes:

```bash
migrate -source file://db/migrations -database "$DATABASE_DSN" up
```

You can run the helper script that ships with the repo:

```bash
sh scripts/migrate.sh
```

## Local development

```bash
# Compile once
make build

# Run with hot reload (rebuilds on code change)
make run

# Run the compiled binary plus dependency checks
make deploy-local   # executes scripts/essential.sh → wait-for-deps → migrate → ./bin/sd-svc-auth
```

The Makefile also exposes:

| Target                                | Description                                   |
| ------------------------------------- | --------------------------------------------- |
| `make test`                           | Run Go tests in `./tests/...`                 |
| `make docs`                           | Serve mdBook docs on port 3000                |
| `make docker-build`                   | Build the image `shinoda4/sd-svc-auth:latest` |
| `make docker-up` / `make docker-down` | Start/stop Docker Compose stack               |

## Running the transports

After `make run` (or executing `go run ./cmd/server`):

- **gRPC** is available on `localhost:$GRPC_PORT` (default `50051`).
- **HTTP gateway** is available on `http://localhost:$HTTP_PORT` (default `8080`).

Health probe:

```bash
grpcurl -plaintext localhost:50051 auth.v1.AuthService.HealthCheck
```

## Docker workflow

```bash
make docker-build
docker run \
  -p 8080:8080 \
  -p 50051:50051 \
  --env-file .env \
  sd-svc-auth:latest
```

### Docker Compose

`deployments/docker-compose.yml` spins up PostgreSQL, Redis, and the service:

```bash
docker-compose -f deployments/docker-compose.yml up -d
docker-compose -f deployments/docker-compose.yml logs -f sd-svc-auth
docker-compose -f deployments/docker-compose.yml down
```

Compose ships with helper scripts:

- `scripts/wait-for-deps.sh` – waits for PostgreSQL/Redis readiness.
- `scripts/run.sh` – runs dependency checks, migrations, then executes the binary inside the container.

## Kubernetes outline

Use the Docker image and mount environment variables via secrets:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sd-svc-auth
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sd-svc-auth
  template:
    metadata:
      labels:
        app: sd-svc-auth
    spec:
      containers:
        - name: auth
          image: shinoda4/sd-svc-auth:latest
          ports:
            - containerPort: 8080
              name: http
            - containerPort: 50051
              name: grpc
          envFrom:
            - secretRef:
                name: sd-svc-auth-env
```

Expose `8080` through an Ingress/Service for HTTP clients and `50051` through a ClusterIP or headless service for gRPC peers. Add a gRPC health probe or an HTTP `/health` handler if required by your platform.

## Troubleshooting

| Symptom                                                     | Checks                                                                                                             |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| Process exits with `missing required environment variables` | `env \| grep -E 'DATABASE_DSN\|REDIS_ADDR\|HTTP_PORT'`                                                             |
| gRPC cannot bind                                            | Ensure `$GRPC_PORT` is free (`lsof -i :50051`).                                                                    |
| Login returns `invalid password`                            | Confirm bcrypt hashes via `psql` and ensure `email_verified` is true.                                              |
| Refresh token fails                                         | Verify Redis is reachable and contains `token:<userID>`.                                                           |
| Emails are not sent                                         | Confirm `EMAIL_ADDRESS`/`EMAIL_PASSWORD`, network egress, and that Gmail app passwords are enabled if using Gmail. |

## Next steps

- Configure monitoring (logs, metrics, tracing) and add TLS termination suited to your environment.
- Review [Database](./database.md) and [API Reference](./api_reference/index.md) for schema and contract details.
