-- +migrate Up

CREATE TABLE IF NOT EXISTS diva_post_interaction (
    id UUID NOT NULL PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES diva_post(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_collection_interaction (
    id UUID NOT NULL PRIMARY KEY,
    collection_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (collection_id) REFERENCES diva_collection(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_message_interaction (
    id UUID NOT NULL PRIMARY KEY,
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (message_id) REFERENCES diva_message(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_post_share (
    interaction_id UUID NOT NULL PRIMARY KEY,
    share_type VARCHAR(50) NOT NULL,
    caption TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (interaction_id) REFERENCES diva_post_interaction(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_post_comment (
    interaction_id UUID NOT NULL PRIMARY KEY,
    reply_to UUID DEFAULT NULL,
    content TEXT NOT NULL,
    edited_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (interaction_id) REFERENCES diva_post_interaction(id) ON DELETE RESTRICT,
    FOREIGN KEY (reply_to) REFERENCES diva_post_comment(interaction_id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_collection_share (
    interaction_id UUID NOT NULL PRIMARY KEY,
    share_type VARCHAR(50) NOT NULL,
    caption TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (interaction_id) REFERENCES diva_collection_interaction(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diva_collection_comment (
    interaction_id UUID NOT NULL PRIMARY KEY,
    reply_to UUID DEFAULT NULL,
    content TEXT NOT NULL,
    edited_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (interaction_id) REFERENCES diva_collection_interaction(id) ON DELETE RESTRICT,
    FOREIGN KEY (reply_to) REFERENCES diva_collection_comment(interaction_id) ON DELETE RESTRICT
);

-- +migrate Down

DROP TABLE IF EXISTS diva_collection_comment;
DROP TABLE IF EXISTS diva_collection_share;
DROP TABLE IF EXISTS diva_post_comment;
DROP TABLE IF EXISTS diva_post_share;
DROP TABLE IF EXISTS diva_message_interaction;
DROP TABLE IF EXISTS diva_collection_interaction;
DROP TABLE IF EXISTS diva_post_interaction;
