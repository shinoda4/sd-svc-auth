# Database

The service uses PostgreSQL as its primary data store.

## Schema

The database consists of a single `users` table.

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

## Migrations

Database migrations are managed using `golang-migrate`.

### Prerequisites

- `golang-migrate` CLI tool.

### Running Migrations

```bash
migrate -source file://db/migrations -database "$DATABASE_DSN" up
```

### Creating a New Migration

```bash
migrate create -ext sql -dir db/migrations -seq [migration_name]
```
