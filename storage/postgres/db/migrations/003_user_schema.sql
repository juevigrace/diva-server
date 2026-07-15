-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_user (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone_number VARCHAR(30) NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL,
    role role_type NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS diva_user_profile(
    user_id UUID NOT NULL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    birth_date TIMESTAMPTZ DEFAULT NULL,
    alias VARCHAR(255) NOT NULL DEFAULT '',
    bio VARCHAR(255) NOT NULL DEFAULT '',
    avatar TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_user_state (
    user_id UUID NOT NULL PRIMARY KEY,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    status user_status_type NOT NULL DEFAULT 'ACTIVE',
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_user_permissions (
    permission_id UUID NOT NULL,
    user_id UUID NOT NULL,
    granted_by UUID DEFAULT NULL,
    granted BOOLEAN NOT NULL DEFAULT FALSE,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ DEFAULT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (permission_id, user_id),
    FOREIGN KEY (permission_id) REFERENCES diva_permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES diva_user(id),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_user_preferences (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    device VARCHAR(100) NOT NULL,
    theme theme_type NOT NULL DEFAULT 'SYSTEM',
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    last_sync_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
