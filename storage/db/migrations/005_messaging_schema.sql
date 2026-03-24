-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diva_chat (
    id UUID NOT NULL PRIMARY KEY,
    created_by UUID NOT NULL,
    chat_type chat_type_type NOT NULL,
    name VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    avatar UUID DEFAULT NULL,
    last_message_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (created_by) REFERENCES diva_user(id) ON DELETE RESTRICT,
    FOREIGN KEY (avatar) REFERENCES diva_media(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_chat_participant (
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role participant_role_type NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_read_at TIMESTAMPTZ DEFAULT NULL,
    added_by UUID NOT NULL,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES diva_chat(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES diva_user(id) ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES diva_user(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_message (
    id UUID NOT NULL PRIMARY KEY,
    chat_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    message_type message_type_type NOT NULL,
    reply_to_id UUID DEFAULT NULL,
    edited_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (chat_id) REFERENCES diva_chat(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES diva_user(id) ON DELETE RESTRICT,
    FOREIGN KEY (reply_to_id) REFERENCES diva_message(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS diva_message_media (
    message_id UUID NOT NULL,
    media_id UUID NOT NULL,
    PRIMARY KEY (message_id, media_id),
    FOREIGN KEY (message_id) REFERENCES diva_message(id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES diva_media(id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diva_message_media;
DROP TABLE IF EXISTS diva_message;
DROP TABLE IF EXISTS diva_chat_participant;
DROP TABLE IF EXISTS diva_chat;
-- +goose StatementEnd
