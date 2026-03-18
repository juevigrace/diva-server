-- name: CreateMessage :exec
INSERT INTO diva_message (id, chat_id, sender_id, content, message_type, reply_to_id)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateMessage :exec
UPDATE diva_message 
SET content = $1, edited_at = NOW(), updated_at = NOW() 
WHERE id = $2 AND deleted_at IS NULL;

-- name: DeleteMessage :exec
UPDATE diva_message 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMessageByID :one
SELECT id, chat_id AS chat_id, sender_id AS sender_id, content, message_type AS message_type, reply_to_id AS reply_to_id, created_at AS created_at, updated_at AS updated_at, deleted_at AS deleted_at, edited_at AS edited_at FROM diva_message 
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMessages :many
SELECT id, chat_id AS chat_id, sender_id AS sender_id, content, message_type AS message_type, reply_to_id AS reply_to_id, created_at AS created_at, updated_at AS updated_at, deleted_at AS deleted_at, edited_at AS edited_at FROM diva_message 
WHERE deleted_at IS NULL 
ORDER BY created_at DESC;

-- name: GetMessagesByChat :many
SELECT id, chat_id AS chat_id, sender_id AS sender_id, content, message_type AS message_type, reply_to_id AS reply_to_id, created_at AS created_at, updated_at AS updated_at, deleted_at AS deleted_at, edited_at AS edited_at FROM diva_message 
WHERE chat_id = $1 AND deleted_at IS NULL 
ORDER BY created_at ASC;

-- name: GetMessagesBySender :many
SELECT id, chat_id AS chat_id, sender_id AS sender_id, content, message_type AS message_type, reply_to_id AS reply_to_id, created_at AS created_at, updated_at AS updated_at, deleted_at AS deleted_at, edited_at AS edited_at FROM diva_message 
WHERE sender_id = $1 AND deleted_at IS NULL 
ORDER BY created_at DESC;
