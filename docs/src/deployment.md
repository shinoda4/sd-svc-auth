# Deployment

## Prerequisites

- **Go**: Version 1.21 or higher.
- **PostgreSQL**: Version 13 or higher.
- **Redis**: Version 6 or higher.

## Environment Variables

The application relies on environment variables for configuration. You can set them directly in your shell or use a `.env` file (if supported by your runner, though this project uses standard `os.Getenv`).

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `SERVER_PORT` | Port to listen on | Yes | `8080` |
| `DATABASE_DSN` | PostgreSQL connection string | Yes | `postgresql://user:pass@localhost:5432/db?sslmode=disable` |
| `REDIS_ADDR` | Redis address | Yes | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | No | `secret` |
| `JWT_SECRET` | Secret key for signing tokens | Yes | `your-256-bit-secret` |
| `EMAIL_ADDRESS` | SMTP email address | Yes | `noreply@example.com` |
| `EMAIL_PASSWORD` | SMTP email password | Yes | `smtp-password` |

## Running the Application

### Local Development

1. **Start Dependencies**: Ensure PostgreSQL and Redis are running.
2. **Set Environment**: Export the required environment variables.
3. **Run**:
   ```bash
   go run cmd/server/main.go
   ```

### Docker

(Assuming a Dockerfile exists)

```bash
docker build -t sd-svc-auth .
docker run -p 8080:8080 --env-file .env sd-svc-auth
```

## Database Migrations

Currently, the application expects the database schema to be set up. Ensure the `users` table exists.

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
