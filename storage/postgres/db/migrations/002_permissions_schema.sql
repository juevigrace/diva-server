-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_permissions (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    action VARCHAR(255) NOT NULL UNIQUE,
    role_level role_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_permissions;
-- +goose StatementEnd
