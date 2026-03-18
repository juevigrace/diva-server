-- name: CreateMessageMedia :exec
INSERT INTO diva_message_media (message_id, media_id)
VALUES ($1, $2);

-- name: DeleteMessageMedia :exec
DELETE FROM diva_message_media 
WHERE message_id = $1 AND media_id = $2;
