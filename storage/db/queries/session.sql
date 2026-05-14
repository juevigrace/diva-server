-- name: CreateSession :exec
INSERT INTO diva_session (id, user_id, access_token, refresh_token, device, status, type, ip_address, user_agent, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdateSession :exec
UPDATE diva_session
SET access_token = $1, refresh_token = $2, device = $3, status = $4, type = $5, ip_address = $6, user_agent = $7, expires_at = $8, updated_at = NOW()
WHERE id = $9;

-- name: UpdateSessionStatus :exec
UPDATE diva_session 
SET  status = $1, updated_at = NOW() 
WHERE id = $2;

-- name: DeleteSession :exec
delete from diva_session
where id = $1
;

-- name: GetSessionByID :one
select
    s.id,
    s.user_id,
    s.access_token,
    s.refresh_token,
    s.device,
    s.status,
    s.type,
    s.ip_address,
    s.user_agent,
    s.expires_at,
    s.created_at,
    s.updated_at
from diva_session s
where s.id = $1
;

-- name: GetSessionsByUser :many
select
    s.id,
    s.user_id,
    s.access_token,
    s.refresh_token,
    s.device,
    s.status,
    s.type,
    s.ip_address,
    s.user_agent,
    s.expires_at,
    s.created_at,
    s.updated_at
from diva_session s
where s.user_id = $1
;
