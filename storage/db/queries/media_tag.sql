-- name: CreateMediaTag :exec
INSERT INTO diva_media_tag (media_id, tag_id)
VALUES ($1, $2);

-- name: UpdateMediaTag :exec
UPDATE diva_media_tag 
SET media_id = $1, tag_id = $2 
WHERE media_id = $3 AND tag_id = $4;

-- name: DeleteMediaTag :exec
DELETE FROM diva_media_tag 
WHERE media_id = $1 AND tag_id = $2;
