-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_type AS ENUM ('USER', 'MODERATOR', 'ADMIN');
CREATE TYPE session_type AS ENUM ('NORMAL', 'TEMPORAL');
CREATE TYPE session_status_type AS ENUM ('ACTIVE', 'EXPIRED', 'CLOSED');
CREATE TYPE theme_type AS ENUM ('LIGHT', 'DARK', 'SYSTEM');
CREATE TYPE media_type_type AS ENUM ('AUDIO', 'IMAGE', 'VIDEO', 'UNSPECIFIED');
CREATE TYPE user_status_type AS ENUM ('ACTIVE', 'SUSPENDED', 'INACTIVE');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS user_status_type;
DROP TYPE IF EXISTS role_type;
DROP TYPE IF EXISTS session_type;
DROP TYPE IF EXISTS session_status_type;
DROP TYPE IF EXISTS theme_type;
DROP TYPE IF EXISTS media_type_type;
-- +goose StatementEnd
