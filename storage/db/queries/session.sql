-- name: GetSessionByID :one
select
    s.id as id,
    s.user_id as userId,
    s.access_token as accessToken,
    s.refresh_token as refreshToken,
    s.device,
    s.status,
    s.type,
    s.ip_address as ipAddress,
    s.user_agent as userAgent,
    s.expires_at as expiresAt,
    s.created_at as createdAt,
    s.updated_at as updatedAt
from diva_session s
where s.id = $1
;

-- name: ListSessionsByUser :many
select
    s.id as id,
    s.user_id as userId,
    s.access_token as accessToken,
    s.refresh_token as refreshToken,
    s.device,
    s.status,
    s.type,
    s.ip_address as ipAddress,
    s.user_agent as userAgent,
    s.expires_at as expiresAt,
    s.created_at as createdAt,
    s.updated_at as updatedAt
from diva_session s
where s.user_id = $1
order by created_at desc
;

-- name: GetSessionByAccessToken :one
select
    s.id as id,
    s.user_id as userId,
    s.access_token as accessToken,
    s.refresh_token as refreshToken,
    s.device,
    s.status,
    s.type,
    s.ip_address as ipAddress,
    s.user_agent as userAgent,
    s.expires_at as expiresAt,
    s.created_at as createdAt,
    s.updated_at as updatedAt
from diva_session s
where s.access_token = $1
;

-- name: GetSessionByRefreshToken :one
select
    s.id as id,
    s.user_id as userId,
    s.access_token as accessToken,
    s.refresh_token as refreshToken,
    s.device,
    s.status,
    s.type,
    s.ip_address as ipAddress,
    s.user_agent as userAgent,
    s.expires_at as expiresAt,
    s.created_at as createdAt,
    s.updated_at as updatedAt
from diva_session s
where s.refresh_token = $1
;

-- name: CreateSession :exec
insert into diva_session (
    id,
    user_id,
    access_token,
    refresh_token,
    device,
    status,
    type,
    ip_address,
    user_agent,
    expires_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
);

-- name: UpdateSession :exec
update diva_session
set
    access_token = $1,
    refresh_token = $2,
    device = $3,
    status = $4,
    type = $5,
    ip_address = $6,
    user_agent = $7,
    expires_at = $8,
    updated_at = now()
where id = $9;

-- name: UpdateSessionStatus :exec
update diva_session
set
    status = $1,
    updated_at = now()
where id = $2;

-- name: DeleteSession :exec
delete from diva_session
where id = $1;

-- name: DeleteSessionsByUser :exec
delete from diva_session
where user_id = $1;

-- name: DeleteExpiredSessions :exec
delete from diva_session
where expires_at < now()
;