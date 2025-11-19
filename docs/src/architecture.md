# Architecture

This project follows the **Standard Go Project Layout** and implements a **Clean Architecture** pattern to ensure separation of concerns, testability, and maintainability.

## Directory Structure

```
.
├── cmd/
│   └── server/         # Application entry point
├── configs/            # Configuration loading logic
├── docs/               # Documentation (mdbook)
├── internal/           # Private application code
│   ├── dto/            # Data Transfer Objects (API request/response structs)
│   ├── transport/      # Transport layer
│   │   └── http/       # HTTP Handlers (Controller layer)
│   ├── model/          # Domain models
│   ├── repo/           # Data Access Layer (Repository pattern)
│   └── service/        # Business Logic Layer
├── pkg/                # Public library code (utilities)
│   ├── email/          # Email sending utility
│   ├── logger/         # Logging utility
│   └── token/          # JWT generation and validation
└── tests/              # Integration and unit tests
```

## Design Patterns

### Repository Pattern
The data access layer (`internal/repo`) abstracts the underlying database and cache technologies. This allows the business logic to depend on interfaces rather than concrete implementations, making it easier to mock data sources for testing or switch databases in the future.

### Service Layer
The service layer (`internal/service`) contains the core business logic. It orchestrates data flow between the handlers and the repositories. It is responsible for validation, transaction management, and enforcing business rules.

### DTO (Data Transfer Object)
We use DTOs (`internal/dto`) to define the structure of data sent to and received from the API. This decouples the internal domain models from the external API representation, allowing them to evolve independently.

### Dependency Injection
Dependencies (like database connections, cache clients, and services) are injected into the components that need them (e.g., Handlers). This promotes loose coupling and makes unit testing straightforward.
