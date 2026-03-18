-- name: CountTag :one
SELECT COUNT(*)
FROM diva_tag
WHERE deleted_at IS NULL;

-- name: CreateTag :exec
INSERT INTO diva_tag (id, tag_name)
VALUES ($1, $2);

-- name: UpdateTag :exec
UPDATE diva_tag 
SET tag_name = $1, updated_at = NOW() 
WHERE id = $2;

-- name: DeleteTag :exec
UPDATE diva_tag
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTags :many
SELECT id, tag_name AS tag_name, created_at AS created_at, updated_at AS updated_at, deleted_at AS deleted_at FROM diva_tag
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTagByID :one
SELECT id, tag_name AS tag_name, created_at AS created_at, updated_at AS updated_at, deleted_at AS deleted_at FROM diva_tag 
WHERE id = $1 AND deleted_at IS NULL;
