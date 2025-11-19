# Architecture

This project follows the **Standard Go Project Layout** and implements a **Clean Architecture** pattern to ensure separation of concerns, testability, and maintainability.

## Directory Structure

```
sd-svc-auth/
├── cmd/
│   └── server/           # Application entry point
├── internal/             # Private application code
│   ├── config/           # Configuration loading and management
│   ├── model/            # Domain models (User entity)
│   ├── repo/             # Data Access Layer (Repository pattern)
│   │   ├── user.go       # User repository (PostgreSQL)
│   │   └── redis.go      # Redis cache operations
│   ├── service/          # Business Logic Layer
│   │   └── auth/         # Authentication service
│   └── transport/        # Transport layer
│       └── grpc/         # gRPC server and HTTP gateway
│           ├── server.go # gRPC server setup and interceptors
│           ├── auth.go   # Authentication endpoints
│           └── jwt.go    # JWT token endpoints
├── pkg/                  # Public library code (utilities)
│   ├── email/            # Email sending utility
│   ├── logger/           # Logging utility
│   └── token/            # JWT generation and validation
├── db/                   # Database migrations
├── docs/                 # Documentation (mdbook)
└── deployments/          # Docker and deployment configs
```

## Transport Layer

The service uses a **gRPC-first architecture** with automatic HTTP/JSON API exposure:

### gRPC Server

- **Protocol**: Native gRPC using Protocol Buffers
- **Service Definition**: `auth.v1.AuthService` (from `sd-grpc-proto`)
- **Port**: Configurable via `GRPC_PORT` (default: 50051)
- **Features**:
  - Unary interceptors for logging and authentication
  - Metadata-based authentication (Bearer tokens)
  - Structured error handling with gRPC status codes

### HTTP Gateway

- **Implementation**: [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- **Port**: Configurable via `HTTP_PORT` (default: 8080)
- **Features**:
  - Automatic JSON ↔ Protobuf conversion
  - RESTful-style HTTP endpoints
  - Standard HTTP authentication (Authorization header)

### Authentication Interceptor

The gRPC server includes a custom authentication interceptor that:

1. **Whitelists** public endpoints (Register, Login, VerifyEmail, ForgotPassword, ResetPassword)
2. **Validates** JWT tokens from metadata for protected endpoints
3. **Injects** claims into the request context for downstream handlers
4. **Returns** appropriate gRPC error codes for authentication failures

## Design Patterns

### Repository Pattern

The data access layer (`internal/repo`) abstracts the underlying database and cache technologies. This allows the business logic to depend on interfaces rather than concrete implementations, making it easier to mock data sources for testing or switch databases in the future.

**Repositories**:
- `UserRepo`: PostgreSQL-based user data persistence
- `RedisCache`: Token blacklist and caching operations

### Service Layer

The service layer (`internal/service`) contains the core business logic. It orchestrates data flow between the transport handlers and the repositories. It is responsible for:

- Input validation
- Business rule enforcement
- Transaction management
- Coordinating between multiple repositories
- Token generation and validation

### Dependency Injection

Dependencies (like database connections, cache clients, and services) are injected into the components that need them. This promotes loose coupling and makes unit testing straightforward.

**Dependency Flow**:
```
main.go
  ↓
  ├─→ UserRepo (PostgreSQL)
  ├─→ RedisCache
  ↓
AuthService (injected with repos)
  ↓
gRPC Server (injected with service)
```

## Data Flow

### Authentication Flow (Login)

1. **Client** sends credentials via gRPC or HTTP
2. **Transport Layer** (gRPC handler) receives request
3. **Service Layer** validates credentials and generates tokens
4. **Repository Layer** queries database and updates cache
5. **Response** returns access and refresh tokens to client

### Protected Endpoint Flow

1. **Client** sends request with Bearer token
2. **Interceptor** extracts and validates token
3. **Claims** injected into context
4. **Handler** processes request with user context
5. **Response** returned to client

## Error Handling

The service uses gRPC status codes for consistent error handling:

- `codes.InvalidArgument`: Invalid input data
- `codes.Unauthenticated`: Missing or invalid authentication
- `codes.AlreadyExists`: Duplicate user registration
- `codes.NotFound`: User not found
- `codes.Internal`: Server errors

These are automatically converted to appropriate HTTP status codes by grpc-gateway.

