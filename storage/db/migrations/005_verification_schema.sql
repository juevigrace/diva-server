-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_action(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    UNIQUE(user_id, name)
);

CREATE TABLE IF NOT EXISTS diva_action_verification (
    action_id UUID NOT NULL PRIMARY KEY,
    token CHAR(6) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    UNIQUE(action_id, token),
    FOREIGN KEY(action_id) REFERENCES diva_user_action(id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_action;
DROP TABLE IF EXISTS diva_action_verification;
-- +goose StatementEnd
