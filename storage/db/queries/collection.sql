-- name: CountCollection :one
SELECT COUNT(*)
FROM diva_collection
WHERE deleted_at IS NULL;

-- name: CreateCollection :exec
INSERT INTO diva_collection (id, owner, cover_media_id, name, description, collection_type, visibility)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateCollection :exec
UPDATE diva_collection 
SET cover_media_id = $1, name = $2, description = $3, visibility = $4, updated_at = NOW()
WHERE id = $5 AND deleted_at IS NULL;

-- name: DeleteCollection :exec
UPDATE diva_collection 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetCollectionByID :one
SELECT id, owner, cover_media_id, name, description, collection_type, visibility, created_at, updated_at, deleted_at FROM diva_collection 
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCollections :many
SELECT id, owner, cover_media_id, name, description, collection_type, visibility, created_at, updated_at, deleted_at FROM diva_collection 
WHERE deleted_at IS NULL 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
