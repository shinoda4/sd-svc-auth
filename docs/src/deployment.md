# Deployment

## Prerequisites

- **Go**: Version 1.21 or higher
- **PostgreSQL**: Version 13 or higher
- **Redis**: Version 6 or higher

## Environment Variables

The application relies on environment variables for configuration. You can set them directly in your shell or use a `.env` file.

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `DATABASE_DSN` | PostgreSQL connection string | Yes | `postgres://user:pass@localhost:5432/db?sslmode=disable` |
| `REDIS_ADDR` | Redis address | Yes | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | No | `secret` |
| `SERVER_HOST` | Server hostname for links | Yes | `localhost` or `https://api.example.com` |
| `SERVER_PORT` | Legacy server port | Yes | `8080` |
| `HTTP_PORT` | HTTP gateway port | Yes | `8080` |
| `GRPC_PORT` | gRPC server port | Yes | `50051` |
| `JWT_SECRET` | Secret key for signing tokens | Yes | `your-256-bit-secret` |
| `JWT_EXPIRE_HOURS` | Token expiration in hours | No | `72` (default) |
| `EMAIL_ADDRESS` | SMTP email address | Yes | `noreply@example.com` |
| `EMAIL_PASSWORD` | SMTP email password | Yes | `smtp-password` |
| `RESET_PASSWORD_URL` | Password reset URL | Yes | `http://localhost:8080/api/v1/reset-password` |

For detailed configuration information, see the [Configuration](./configuration.md) page.

## Running the Application

### Local Development

1. **Start Dependencies**: Ensure PostgreSQL and Redis are running.

2. **Set Environment**: Create a `.env` file based on `.env.example`:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run Database Migrations**:
   ```bash
   migrate -source file://db/migrations -database "$DATABASE_DSN" up
   ```

4. **Run the Service**:
   ```bash
   # Using Make
   make run
   
   # Or directly
   go run cmd/server/main.go
   ```

The service will start:
- **gRPC server** on port `50051` (or `$GRPC_PORT`)
- **HTTP gateway** on port `8080` (or `$HTTP_PORT`)

### Using Make Commands

```bash
# Build binary
make build

# Run service
make run

# Run tests
make test

# Build Docker image
make docker

# Start with Docker Compose
make up

# Initialize database
make init-db
```

### Docker

#### Build and Run

```bash
# Build image
docker build -t sd-svc-auth:latest .

# Run container
docker run -p 8080:8080 -p 50051:50051 \
  -e DATABASE_DSN="postgres://..." \
  -e REDIS_ADDR="redis:6379" \
  -e JWT_SECRET="your-secret" \
  sd-svc-auth:latest
```

#### Docker Compose

```bash
# Start all services
docker-compose -f deployments/docker-compose.yml up -d

# View logs
docker-compose -f deployments/docker-compose.yml logs -f

# Stop services
docker-compose -f deployments/docker-compose.yml down
```

## Database Setup

### Migrations

The service uses `golang-migrate` for database migrations.

#### Install golang-migrate

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

#### Run Migrations

```bash
migrate -source file://db/migrations -database "$DATABASE_DSN" up
```

#### Create New Migration

```bash
migrate create -ext sql -dir db/migrations -seq add_new_field
```

### Manual Schema Setup

If not using migrations, create the schema manually:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    verify_token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Production Deployment

### Security Checklist

- [ ] Use strong `JWT_SECRET` (minimum 32 characters)
- [ ] Enable SSL/TLS for PostgreSQL (`sslmode=require`)
- [ ] Use Redis password authentication
- [ ] Set appropriate `JWT_EXPIRE_HOURS` (24-72 hours recommended)
- [ ] Use app-specific passwords for email
- [ ] Configure proper firewall rules
- [ ] Enable HTTPS for HTTP gateway
- [ ] Use TLS for gRPC in production

### Environment Configuration

Create a production `.env` file with secure values:

```bash
DATABASE_DSN=postgres://auth_user:strong_password@db.example.com:5432/auth_prod?sslmode=require
REDIS_ADDR=redis.example.com:6379
REDIS_PASSWORD=redis_strong_password
SERVER_HOST=https://api.example.com
HTTP_PORT=8080
GRPC_PORT=50051
JWT_SECRET=<generated-with-openssl-rand-base64-32>
JWT_EXPIRE_HOURS=24
EMAIL_ADDRESS=noreply@example.com
EMAIL_PASSWORD=smtp_app_password
RESET_PASSWORD_URL=https://app.example.com/reset-password
```

### Kubernetes Deployment

Example Kubernetes deployment:

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
        image: sd-svc-auth:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 50051
          name: grpc
        env:
        - name: DATABASE_DSN
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: database-dsn
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secrets
              key: jwt-secret
        # ... other env vars
```

## Health Checks

The service doesn't currently expose health check endpoints. Consider adding:

- `/health` for HTTP gateway
- gRPC health check service

## Monitoring

Recommended monitoring:

- **Metrics**: Prometheus metrics for request counts, latency, errors
- **Logging**: Structured logging to stdout (captured by container runtime)
- **Tracing**: OpenTelemetry for distributed tracing

## Troubleshooting

### Service Won't Start

1. Check database connectivity:
   ```bash
   psql "$DATABASE_DSN"
   ```

2. Check Redis connectivity:
   ```bash
   redis-cli -h <host> -p <port> ping
   ```

3. Verify environment variables are set:
   ```bash
   env | grep -E 'DATABASE_DSN|REDIS_ADDR|JWT_SECRET'
   ```

### Database Connection Errors

- Ensure PostgreSQL is running and accessible
- Verify connection string format
- Check firewall rules
- Verify database user permissions

### Redis Connection Errors

- Ensure Redis is running
- Verify Redis address and port
- Check Redis password if authentication is enabled

