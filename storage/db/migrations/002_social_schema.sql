-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_follow (
    user_id UUID NOT NULL,
    followed UUID NOT NULL,
    followed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, followed),
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE,
    FOREIGN KEY (followed) REFERENCES diva_user(id) ON DELETE CASCADE,
    CONSTRAINT no_self_follow CHECK (user_id != followed)
);

CREATE TABLE IF NOT EXISTS diva_post (
    id UUID NOT NULL PRIMARY KEY,
    author_id UUID NOT NULL,
    text TEXT NOT NULL DEFAULT '',
    visibility VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (author_id) REFERENCES diva_user(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_post_attachment (
    id UUID NOT NULL PRIMARY KEY,
    post_id UUID NOT NULL,
    media_id UUID DEFAULT NULL,
    collection_id UUID DEFAULT NULL,
    FOREIGN KEY (post_id) REFERENCES diva_post(id) ON DELETE RESTRICT,
    FOREIGN KEY (media_id) REFERENCES diva_media(id) ON DELETE RESTRICT,
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE RESTRICT,
    CONSTRAINT attachment_type CHECK (
        (media_id IS NOT NULL AND collection_id IS NULL) OR
        (media_id IS NULL AND collection_id IS NOT NULL)
    )
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_post_attachment;
DROP TABLE IF EXISTS diva_post;
DROP TABLE IF EXISTS diva_follow;
-- +goose StatementEnd
