-- name: CountChat :one
SELECT COUNT(*)
FROM diva_chat
WHERE deleted_at IS NULL;

-- name: CreateChat :exec
INSERT INTO diva_chat (id, created_by, chat_type, name, description, avatar)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateChat :exec
UPDATE diva_chat 
SET name = $1, description = $2, avatar = $3, updated_at = NOW()
WHERE id = $4 AND deleted_at IS NULL;

-- name: DeleteChat :exec
UPDATE diva_chat 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetChatByID :one
SELECT 
  diva_chat.id,
  diva_chat.chat_type AS chat_type,
  diva_chat.name,
  diva_chat.description,
  diva_chat.avatar,
  diva_chat.created_by AS created_by,
  diva_chat.created_at AS created_at,
  diva_chat.updated_at AS updated_at,
  diva_chat.deleted_at AS deleted_at
FROM diva_chat
LEFT JOIN diva_chat_participant
ON diva_chat.id = diva_chat_participant.chat_id
LEFT JOIN diva_user
ON diva_chat_participant.user_id = diva_user.id
LEFT JOIN diva_media
ON diva_chat.avatar = diva_media.id
LEFT JOIN diva_media_tag
ON diva_media.id = diva_media_tag.media_id
LEFT JOIN diva_tag
ON diva_media_tag.tag_id = diva_tag.id
WHERE diva_chat.id = $1 AND diva_chat.deleted_at IS NULL;

-- name: ListChats :many
SELECT 
  id,
  chat_type,
  name,
  description,
  avatar,
  created_by,
  created_at,
  updated_at,
  deleted_at
FROM diva_chat
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
