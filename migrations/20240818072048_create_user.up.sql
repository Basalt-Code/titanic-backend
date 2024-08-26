CREATE TABLE IF NOT EXISTS users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nickname            VARCHAR NOT NULL UNIQUE,
    email               TEXT UNIQUE NOT NULL UNIQUE,
    password            TEXT NOT NULL,
    role                VARCHAR(50) NOT NULL,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP
);
