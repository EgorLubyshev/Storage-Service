-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(), -- Требует расширения pg_uuidv7
    name VARCHAR(50) NOT NULL UNIQUE,
    password_hash CHAR(60) NOT NULL,
    rights BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_password_hash_length CHECK (LENGTH(password_hash) = 60)
);

-- Индекс для имени (уже есть благодаря UNIQUE, но можно явно указать)
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
-- Индекс для created_at
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
