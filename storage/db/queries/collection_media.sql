-- name: CreateCollectionMedia :exec
INSERT INTO diva_collection_media (collection_id, media_id, position, added_by, score)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateCollectionMedia :exec
UPDATE diva_collection_media
SET position = $1, score = $2
WHERE collection_id = $3 AND media_id = $4;

-- name: DeleteCollectionMedia :exec
DELETE FROM diva_collection_media
WHERE collection_id = $1 AND media_id = $2;
