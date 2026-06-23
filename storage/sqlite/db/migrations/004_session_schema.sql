-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_session (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    device TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL CHECK (type IN ('NORMAL', 'TEMPORAL')),
    status TEXT NOT NULL CHECK (status IN ('ACTIVE', 'EXPIRED', 'CLOSED')),
    ip_address TEXT NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_session;
-- +goose StatementEnd
