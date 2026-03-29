-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_user (
    id UUID NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    birth_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    phone_number VARCHAR(50) NOT NULL DEFAULT '',
    alias VARCHAR(50) NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    user_verified BOOLEAN NOT NULL DEFAULT FALSE,
    role role_type NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS diva_user_pending_actions(
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    action_name VARCHAR(100) NOT NULL,
    UNIQUE(user_id, action_name),
    FOREIGN KEY(user_id) REFERENCES diva_user(id)
);

CREATE TABLE IF NOT EXISTS diva_user_preferences (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    theme theme_type NOT NULL DEFAULT 'SYSTEM',
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    last_sync_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_session (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    access_token VARCHAR(500) NOT NULL,
    refresh_token VARCHAR(500) NOT NULL,
    device VARCHAR(100) NOT NULL DEFAULT '',
    status session_status_type NOT NULL,
    ip_address VARCHAR(45) NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_permissions (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT '',
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    role_level role_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
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
    FOREIGN KEY (permission_id) REFERENCES diva_permissions(id),
    FOREIGN KEY (granted_by) REFERENCES diva_user(id),
    FOREIGN KEY (user_id) REFERENCES diva_user(id)
);

CREATE TABLE IF NOT EXISTS diva_email_verification_tokens (
    user_id UUID NOT NULL,
    action_id UUID NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, token, action_id),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE,
    FOREIGN KEY (action_id) REFERENCES diva_user_pending_actions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_diva_user_permissions_user_id ON diva_user_permissions (user_id);
CREATE INDEX IF NOT EXISTS idx_diva_user_permissions_permission_id ON diva_user_permissions (permission_id);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_diva_user_permissions_user_id;
DROP INDEX IF EXISTS idx_diva_user_permissions_permission_id;
DROP TABLE IF EXISTS diva_email_verification_tokens;
DROP TABLE IF EXISTS diva_user_permissions;
DROP TABLE IF EXISTS diva_permissions;
DROP TABLE IF EXISTS diva_session;
DROP TABLE IF EXISTS diva_user_preferences;
DROP TABLE IF EXISTS diva_user_pending_actions;
DROP TABLE IF EXISTS diva_user;
-- +goose StatementEnd
