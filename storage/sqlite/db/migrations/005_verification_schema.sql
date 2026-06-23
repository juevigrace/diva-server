-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_action(
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id TEXT NOT NULL,
    UNIQUE(user_id, name),
    FOREIGN KEY(user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_action_verification (
    action_id TEXT NOT NULL PRIMARY KEY,
    token TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP DEFAULT NULL,
    UNIQUE(action_id, token),
    FOREIGN KEY(action_id) REFERENCES diva_action(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_action;
DROP TABLE IF EXISTS diva_action_verification;
-- +goose StatementEnd
