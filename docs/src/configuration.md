# Configuration

This page documents all configuration options for the authentication service.

## Environment Variables

The service is configured entirely through environment variables. Create a `.env` file in the project root or set these variables in your deployment environment.

### Database Configuration

#### `DATABASE_DSN`

**Required**: Yes  
**Description**: PostgreSQL connection string  
**Format**: `postgres://username:password@host:port/database?sslmode=disable`  
**Example**:
```bash
DATABASE_DSN=postgres://sd_auth:sd_pass@localhost:5432/sd_auth?sslmode=disable
```

**Notes**:
- Use `sslmode=require` in production
- Ensure the database exists before starting the service
- The user must have CREATE, READ, UPDATE, DELETE permissions

---

### Redis Configuration

#### `REDIS_ADDR`

**Required**: Yes  
**Description**: Redis server address  
**Format**: `host:port`  
**Example**:
```bash
REDIS_ADDR=localhost:6379
```

#### `REDIS_PASSWORD`

**Required**: No  
**Description**: Redis authentication password  
**Default**: Empty (no authentication)  
**Example**:
```bash
REDIS_PASSWORD=your_redis_password
```

**Notes**:
- Redis is used for token blacklisting and caching
- Leave empty if Redis doesn't require authentication

---

### Server Configuration

#### `SERVER_HOST`

**Required**: Yes  
**Description**: Server hostname for generating verification links  
**Example**:
```bash
SERVER_HOST=localhost
```

**Production Example**:
```bash
SERVER_HOST=https://api.example.com
```

#### `SERVER_PORT`

**Required**: Yes  
**Description**: Legacy server port (for compatibility)  
**Example**:
```bash
SERVER_PORT=8080
```

#### `HTTP_PORT`

**Required**: Yes  
**Description**: HTTP gateway port  
**Default**: `8080`  
**Example**:
```bash
HTTP_PORT=8080
```

#### `GRPC_PORT`

**Required**: Yes  
**Description**: gRPC server port  
**Default**: `50051`  
**Example**:
```bash
GRPC_PORT=50051
```

---

### JWT Configuration

#### `JWT_SECRET`

**Required**: Yes  
**Description**: Secret key for signing JWT tokens  
**Security**: Use a strong, random string (minimum 32 characters)  
**Example**:
```bash
JWT_SECRET=your-super-secret-key-change-this-in-production
```

**Generate a secure secret**:
```bash
openssl rand -base64 32
```

#### `JWT_EXPIRE_HOURS`

**Required**: No  
**Description**: Access token expiration time in hours  
**Default**: `72` (3 days)  
**Example**:
```bash
JWT_EXPIRE_HOURS=24
```

**Notes**:
- Shorter expiration times are more secure but require more frequent refreshes
- Refresh tokens typically have longer expiration (configured in code)

---

### Email Configuration

#### `EMAIL_ADDRESS`

**Required**: Yes  
**Description**: SMTP email address for sending emails  
**Example**:
```bash
EMAIL_ADDRESS=noreply@example.com
```

**Gmail Example**:
```bash
EMAIL_ADDRESS=your-email@gmail.com
```

#### `EMAIL_PASSWORD`

**Required**: Yes  
**Description**: SMTP email password or app-specific password  
**Example**:
```bash
EMAIL_PASSWORD=your_app_specific_password
```

**Notes**:
- For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833)
- Never commit real passwords to version control

---

### Password Reset Configuration

#### `RESET_PASSWORD_URL`

**Required**: Yes  
**Description**: Base URL for password reset links sent in emails  
**Example**:
```bash
RESET_PASSWORD_URL=http://localhost:8080/api/v1/reset-password
```

**Production Example**:
```bash
RESET_PASSWORD_URL=https://app.example.com/reset-password
```

---

## Example Configuration Files

### Development (`.env`)

```bash
# Database
DATABASE_DSN=postgres://sd_auth:sd_pass@localhost:5432/sd_auth?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# Server
SERVER_HOST=localhost
SERVER_PORT=8080
HTTP_PORT=8080
GRPC_PORT=50051

# JWT
JWT_SECRET=dev-secret-change-in-production
JWT_EXPIRE_HOURS=72

# Email (Gmail example)
EMAIL_ADDRESS=your-email@gmail.com
EMAIL_PASSWORD=your-app-password

# Password Reset
RESET_PASSWORD_URL=http://localhost:8080/api/v1/reset-password
```

### Production (`.env.production`)

```bash
# Database
DATABASE_DSN=postgres://auth_user:strong_password@db.example.com:5432/auth_prod?sslmode=require

# Redis
REDIS_ADDR=redis.example.com:6379
REDIS_PASSWORD=redis_strong_password

# Server
SERVER_HOST=https://api.example.com
SERVER_PORT=8080
HTTP_PORT=8080
GRPC_PORT=50051

# JWT
JWT_SECRET=<generated-with-openssl-rand-base64-32>
JWT_EXPIRE_HOURS=24

# Email
EMAIL_ADDRESS=noreply@example.com
EMAIL_PASSWORD=smtp_app_password

# Password Reset
RESET_PASSWORD_URL=https://app.example.com/reset-password
```

## Docker Configuration

When using Docker, pass environment variables via:

### Docker Run

```bash
docker run -p 8080:8080 -p 50051:50051 \
  -e DATABASE_DSN="postgres://..." \
  -e REDIS_ADDR="redis:6379" \
  -e JWT_SECRET="your-secret" \
  sd-svc-auth:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  auth:
    image: sd-svc-auth:latest
    ports:
      - "8080:8080"
      - "50051:50051"
    environment:
      DATABASE_DSN: postgres://sd_auth:sd_pass@postgres:5432/sd_auth?sslmode=disable
      REDIS_ADDR: redis:6379
      JWT_SECRET: your-secret-key
      # ... other variables
    depends_on:
      - postgres
      - redis
```

## Security Best Practices

1. **Never commit `.env` files** with real credentials to version control
2. **Use strong secrets** for `JWT_SECRET` (minimum 32 characters)
3. **Enable SSL/TLS** in production (`sslmode=require` for PostgreSQL)
4. **Use app-specific passwords** for email services
5. **Rotate secrets regularly** in production environments
6. **Use environment-specific configurations** (dev, staging, production)
7. **Restrict database user permissions** to only what's needed
