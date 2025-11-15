CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    username      VARCHAR NOT NULL UNIQUE,
    email         TEXT    NOT NULL UNIQUE,
    password_hash TEXT    NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now()
);
