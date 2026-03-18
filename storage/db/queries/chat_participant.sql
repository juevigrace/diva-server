-- name: CreateChatParticipant :exec
INSERT INTO diva_chat_participant (chat_id, user_id, role, added_by)
VALUES ($1, $2, $3, $4);

-- name: UpdateChatParticipant :exec
UPDATE diva_chat_participant 
SET role = $1, last_read_at = $2
WHERE chat_id = $3 AND user_id = $4;

-- name: DeleteChatParticipant :exec
DELETE FROM diva_chat_participant 
WHERE chat_id = $1 AND user_id = $2;
