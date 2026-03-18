-- name: CountPost :one
SELECT COUNT(*)
FROM diva_post
WHERE deleted_at IS NULL;

-- name: CreatePost :exec
INSERT INTO diva_post (id, author_id, text, visibility)
VALUES ($1, $2, $3, $4);

-- name: DeletePost :exec
UPDATE diva_post 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPostByID :one
SELECT 
  id,
  author_id AS author_id,
  text,
  visibility,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_post
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPosts :many
SELECT 
  id,
  author_id AS author_id,
  text,
  visibility,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_post
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
