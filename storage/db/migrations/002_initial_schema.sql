-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_user (
    id UUID NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    birth_date TIMESTAMPTZ NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    alias VARCHAR(50) NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    user_verified BOOLEAN NOT NULL DEFAULT FALSE,
    role role_type NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
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
    granted_at TIMESTAMPTZ DEFAULT NULL,
    expires_at TIMESTAMPTZ DEFAULT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (permission_id, user_id),
    FOREIGN KEY (permission_id) REFERENCES diva_permissions(id),
    FOREIGN KEY (granted_by) REFERENCES diva_user(id),
    FOREIGN KEY (user_id) REFERENCES diva_user(id)
);

CREATE TABLE IF NOT EXISTS diva_email_verification_tokens (
    user_id UUID NOT NULL PRIMARY KEY,
    token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_media (
    id UUID NOT NULL PRIMARY KEY,
    submitted_by UUID NOT NULL,
    url TEXT NOT NULL,
    alt_text TEXT NOT NULL DEFAULT '',
    media_type media_type_type NOT NULL,
    file_size BIGINT NOT NULL,
    width INTEGER NOT NULL DEFAULT 0,
    height INTEGER NOT NULL DEFAULT 0,
    duration INTEGER NOT NULL DEFAULT 0,
    visibility visibility_type NOT NULL,
    sensitive_content BOOLEAN NOT NULL DEFAULT FALSE,
    adult_content BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (submitted_by) REFERENCES diva_user(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_tag (
    id UUID NOT NULL PRIMARY KEY,
    tag_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS diva_media_tag (
    media_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    PRIMARY KEY (media_id, tag_id),
    FOREIGN KEY (media_id) REFERENCES diva_media(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES diva_tag(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_collection (
    id UUID NOT NULL PRIMARY KEY,
    owner UUID NOT NULL,
    cover_media_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    collection_type collection_type_type NOT NULL,
    visibility visibility_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (owner) REFERENCES diva_user(id) ON DELETE RESTRICT,
    FOREIGN KEY (cover_media_id) REFERENCES diva_media(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_playlist_metadata (
    collection_id UUID NOT NULL PRIMARY KEY,
    is_collaborative BOOLEAN NOT NULL DEFAULT FALSE,
    allow_suggestions BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_playlist_contributor (
    collection_id UUID NOT NULL,
    contributor_id UUID NOT NULL,
    FOREIGN KEY (contributor_id) REFERENCES diva_user(id),
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, contributor_id)
);

CREATE TABLE IF NOT EXISTS diva_playlist_suggestions (
    collection_id UUID NOT NULL,
    suggester_id UUID NOT NULL,
    media_id UUID NOT NULL,
    status moderation_status_type NOT NULL,
    suggested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (media_id) REFERENCES diva_media(id),
    FOREIGN KEY (suggester_id) REFERENCES diva_user(id),
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, suggester_id, media_id)
);

CREATE TABLE IF NOT EXISTS diva_mix_metadata (
    collection_id UUID NOT NULL PRIMARY KEY,
    algorithm_type VARCHAR(50) NOT NULL DEFAULT 'trending',
    time_window_hours INTEGER NOT NULL DEFAULT 24,
    content_weight DECIMAL(3,2) NOT NULL DEFAULT 0.7,
    freshness_weight DECIMAL(3,2) NOT NULL DEFAULT 0.3,
    min_engagement_score INTEGER NOT NULL DEFAULT 10,
    excluded_tags TEXT[] NOT NULL DEFAULT '{}',
    auto_refresh BOOLEAN NOT NULL DEFAULT TRUE,
    refresh_interval_seconds INTEGER NOT NULL DEFAULT 3600,
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_collection_media (
    collection_id UUID NOT NULL,
    media_id UUID NOT NULL,
    position INTEGER NOT NULL,
    added_by UUID NOT NULL,
    score DECIMAL(5,4) NOT NULL DEFAULT 0.0000,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES diva_media(id) ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES diva_user(id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, media_id)
);

CREATE INDEX IF NOT EXISTS idx_diva_user_permissions_user_id ON diva_user_permissions (user_id);
CREATE INDEX IF NOT EXISTS idx_diva_user_permissions_permission_id ON diva_user_permissions (permission_id);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_diva_user_permissions_user_id;
DROP INDEX IF EXISTS idx_diva_user_permissions_permission_id;
DROP TABLE IF EXISTS diva_collection_media;
DROP TABLE IF EXISTS diva_mix_metadata;
DROP TABLE IF EXISTS diva_playlist_suggestions;
DROP TABLE IF EXISTS diva_playlist_contributor;
DROP TABLE IF EXISTS diva_playlist_metadata;
DROP TABLE IF EXISTS diva_collection;
DROP TABLE IF EXISTS diva_media_tag;
DROP TABLE IF EXISTS diva_tag;
DROP TABLE IF EXISTS diva_media;
DROP TABLE IF EXISTS diva_email_verification_tokens;
DROP TABLE IF EXISTS diva_user_permissions;
DROP TABLE IF EXISTS diva_permissions;
DROP TABLE IF EXISTS diva_session;
DROP TABLE IF EXISTS diva_user_preferences;
DROP TABLE IF EXISTS diva_user;
-- +goose StatementEnd
