-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'USER' CHECK (role IN ('USER', 'MODERATOR', 'ADMIN')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS diva_user_profile(
    user_id TEXT NOT NULL PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    birth_date TIMESTAMP DEFAULT NULL,
    alias TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    avatar TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_user_state (
    user_id TEXT NOT NULL PRIMARY KEY,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SUSPENDED', 'INACTIVE')),
    last_active_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_user_permissions (
    permission_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    granted_by TEXT DEFAULT NULL,
    granted BOOLEAN NOT NULL DEFAULT FALSE,
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (permission_id, user_id),
    FOREIGN KEY (permission_id) REFERENCES diva_permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES diva_user(id),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_user_preferences (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    device TEXT NOT NULL,
    theme TEXT NOT NULL DEFAULT 'SYSTEM' CHECK (theme IN ('LIGHT', 'DARK', 'SYSTEM')),
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    language TEXT NOT NULL DEFAULT 'en',
    last_sync_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_user_state;
DROP TABLE IF EXISTS diva_user;
DROP TABLE IF EXISTS diva_user_profile;
DROP TABLE IF EXISTS diva_user_permissions;
DROP TABLE IF EXISTS diva_user_preferences;
-- +goose StatementEnd
