-- name: CreateSession :exec
INSERT INTO diva_session (id, user_id, access_token, refresh_token, device, status, ip_address, user_agent, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: UpdateSession :exec
UPDATE diva_session 
SET access_token = $1, refresh_token = $2, device = $3, status = $4, ip_address = $5, user_agent = $6, expires_at = $7, updated_at = NOW() 
WHERE id = $8;

-- name: DeleteSessionByUserID :exec
DELETE FROM diva_session
WHERE user_id = $1;

-- name: DeleteSession :exec
DELETE FROM diva_session
WHERE id = $1;

-- name: GetSessionByID :one
SELECT
  s.id,
  s.access_token,
  s.refresh_token,
  s.device,
  s.status,
  s.ip_address,
  s.user_agent,
  s.expires_at,
  s.created_at,
  s.updated_at,
  u.id AS user_id,
  u.email,
  u.username,
  u.user_verified,
  u.role,
  u.created_at,
  u.updated_at
FROM diva_session s
LEFT JOIN diva_user u ON s.user_id = u.id
WHERE s.id = $1 AND u.deleted_at IS NULL;
