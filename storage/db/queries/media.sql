-- name: CountMedia :one
SELECT COUNT(*)
FROM diva_media
WHERE deleted_at IS NULL;

-- name: CreateMedia :exec
INSERT INTO diva_media (id, submitted_by, url, alt_text, media_type, file_size, width, height, duration, visibility, sensitive_content, adult_content)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: UpdateMedia :exec
UPDATE diva_media 
SET url = $1, alt_text = $2, media_type = $3, file_size = $4, width = $5, height = $6, duration = $7, visibility = $8, sensitive_content = $9, adult_content = $10, updated_at = NOW()
WHERE id = $11 AND deleted_at IS NULL;

-- name: DeleteMedia :exec
UPDATE diva_media 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMediaByID :one
SELECT 
  id,
  submitted_by AS submitted_by,
  url,
  alt_text AS alt_text,
  media_type AS media_type,
  file_size AS file_size,
  width,
  height,
  duration,
  visibility,
  sensitive_content AS sensitive_content,
  adult_content AS adult_content,
  published_at AS published_at,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_media
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMedia :many
SELECT 
  id,
  submitted_by AS submitted_by,
  url,
  alt_text AS alt_text,
  media_type AS media_type,
  file_size AS file_size,
  width,
  height,
  duration,
  visibility,
  sensitive_content AS sensitive_content,
  adult_content AS adult_content,
  published_at AS published_at,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_media
WHERE deleted_at IS NULL 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
