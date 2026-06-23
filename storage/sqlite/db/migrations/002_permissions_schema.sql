-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_permissions (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    action TEXT NOT NULL UNIQUE,
    role_level TEXT NOT NULL CHECK (role_level IN ('USER', 'MODERATOR', 'ADMIN')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_permissions;
-- +goose StatementEnd
