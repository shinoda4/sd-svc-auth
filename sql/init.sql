CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
    );

-- Enable pgcrypto for gen_random_uuid() if not enabled:
CREATE EXTENSION IF NOT EXISTS pgcrypto;
