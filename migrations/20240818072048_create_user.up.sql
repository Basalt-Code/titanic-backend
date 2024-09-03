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

CREATE TABLE IF NOT EXISTS user_sessions (
    user_id             UUID NOT NULL,
    session_id          UUID NOT NULL,
    refresh_token       TEXT NOT NULL,
    expired_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP WITH TIME ZONE
);