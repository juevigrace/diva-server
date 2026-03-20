-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_type AS ENUM ('USER', 'MODERATOR', 'ADMIN');
CREATE TYPE session_status_type AS ENUM ('ACTIVE', 'EXPIRED', 'CLOSED');
CREATE TYPE theme_type AS ENUM ('LIGHT', 'DARK', 'SYSTEM');
CREATE TYPE visibility_type AS ENUM ('PUBLIC', 'PRIVATE', 'FRIENDS', 'UNSPECIFIED');
CREATE TYPE media_type_type AS ENUM ('AUDIO', 'IMAGE', 'VIDEO', 'UNSPECIFIED');
CREATE TYPE collection_type_type AS ENUM ('ALBUM', 'PLAYLIST', 'MIX', 'FAVORITES', 'FEATURED', 'TRENDING', 'UNSPECIFIED');
CREATE TYPE chat_type_type AS ENUM ('DIRECT', 'GROUP', 'UNSPECIFIED');
CREATE TYPE participant_role_type AS ENUM ('OWNER', 'ADMIN', 'MEMBER', 'UNSPECIFIED');
CREATE TYPE message_type_type AS ENUM ('TEXT', 'MEDIA', 'SYSTEM', 'UNSPECIFIED');
CREATE TYPE reaction_type_type AS ENUM ('LIKE', 'COMMENT', 'BOOKMARK', 'SHARE', 'UNSPECIFIED');
CREATE TYPE share_type_type AS ENUM ('DIRECT', 'EXTERNAL', 'EMBED', 'DOWNLOAD', 'UNSPECIFIED');
CREATE TYPE moderation_status_type AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'HIDDEN', 'UNSPECIFIED');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS moderation_status_type;
DROP TYPE IF EXISTS share_type_type;
DROP TYPE IF EXISTS reaction_type_type;
DROP TYPE IF EXISTS message_type_type;
DROP TYPE IF EXISTS participant_role_type;
DROP TYPE IF EXISTS chat_type_type;
DROP TYPE IF EXISTS collection_type_type;
DROP TYPE IF EXISTS media_type_type;
DROP TYPE IF EXISTS visibility_type;
DROP TYPE IF EXISTS theme_type;
DROP TYPE IF EXISTS session_status_type;
DROP TYPE IF EXISTS role_type;
-- +goose StatementEnd
