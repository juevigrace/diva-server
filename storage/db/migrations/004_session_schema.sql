-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_session (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    device VARCHAR(100) NOT NULL DEFAULT '',
    type session_type NOT NULL,
    status session_status_type NOT NULL,
    ip_address VARCHAR(45) NOT NULL DEFAULT '',
    user_agent VARCHAR(255) NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_session;
-- +goose StatementEnd
