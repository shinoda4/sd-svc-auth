# Database

`sd-svc-auth` uses PostgreSQL as its source of truth for users, verification tokens, and password reset metadata. All schema changes live under `db/migrations` and are applied with `golang-migrate`.

## Schema

The service currently relies on a single table:

```sql
CREATE TABLE IF NOT EXISTS users (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username           VARCHAR NOT NULL UNIQUE,
    email              TEXT    NOT NULL UNIQUE,
    password_hash      TEXT    NOT NULL,
    email_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    verify_token       VARCHAR(64),
    reset_token        TEXT,
    reset_token_expire TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Fields map directly to the `internal/model.User` struct and the repository methods:

- `verify_token` – populated when users register; cleared once email is verified.
- `reset_token` / `reset_token_expire` – issued during the password-reset flow and invalidated after success.
- `email_verified` – acts as a guard in `service.Login`.

## Migrations

Migrations are timestamped `.up.sql`/`.down.sql` files:

```
db/migrations/
├── 20251115170148_create_users_table.up.sql
├── 20251115170148_create_users_table.down.sql
├── 20251118073509_add_reset_password_to_users.up.sql
└── 20251118073509_add_reset_password_to_users.down.sql
```

Apply them with:

```bash
migrate -source file://db/migrations -database "$DATABASE_DSN" up
```

Rollback (use with caution):

```bash
migrate -source file://db/migrations -database "$DATABASE_DSN" down 1
```

To create a new migration:

```bash
migrate create -ext sql -dir db/migrations -seq add_new_table
```

## Operational tips

- Enable the `pgcrypto` extension (for `gen_random_uuid()`), e.g. `CREATE EXTENSION IF NOT EXISTS pgcrypto;`.
- Monitor `users` for potential growth. Add indexes on `email` and `verify_token` if you expect large datasets (the migrations already enforce `UNIQUE` constraints on email/username).
- Use connection pooling (pgBouncer or cloud equivalents) when scaling horizontally.
- Regularly vacuum/analyze, especially if password reset tokens are churned frequently.
- Back up the database before rotating `JWT_SECRET` to preserve a restore point in case clients need to re-register.

## Local inspection

Useful queries during debugging:

```sql
-- Find a user by email
SELECT id, email, username, email_verified FROM users WHERE email = 'user@example.com';

-- List users awaiting verification
SELECT email, verify_token FROM users WHERE email_verified = false;

-- Inspect password reset tokens
SELECT email, reset_token, reset_token_expire FROM users WHERE reset_token IS NOT NULL;
```

## Related caches

PostgreSQL holds the source of truth, but Redis augments it:

- `token:<userID>` – refresh token currently issued to the user.
- `blacklist:<token>` – access token blacklist used by `service.Logout`.

Both Redis keys expire automatically (matching the JWT TTL) so the database remains lightweight.
