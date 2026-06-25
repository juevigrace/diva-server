-- name: GetSessionByID :one
select
    s.id,
    s.user_id,
    s.access_token,
    s.refresh_token,
    s.device,
    s.type,
    s.status,
    s.ip_address,
    s.user_agent,
    s.access_expires_at,
    s.refresh_expires_at,
    s.created_at,
    s.updated_at
from diva_session s
where s.id = $1
;

-- name: ListSessionsByUser :many
select
    s.id,
    s.user_id,
    s.access_token,
    s.refresh_token,
    s.device,
    s.type,
    s.status,
    s.ip_address,
    s.user_agent,
    s.access_expires_at,
    s.refresh_expires_at,
    s.created_at,
    s.updated_at
from diva_session s
where s.user_id = $1
order by created_at desc
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
    access_expires_at,
    refresh_expires_at
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
    $10,
    $11
);

-- name: UpdateSession :exec
update diva_session set
    access_token = $1,
    refresh_token = $2,
    ip_address = $3,
    access_expires_at = $4,
    refresh_expires_at = $5,
    updated_at = now()
where id = $6;

-- name: UpdateSessionStatus :exec
update diva_session set
    status = $1,
    updated_at = now()
where id = $2;

-- name: DeleteSession :exec
delete from diva_session
where id = $1
;

-- name: DeleteSessionsByUser :exec
delete from diva_session
where user_id = $1
;

-- name: DeleteExpiredSessions :exec
delete from diva_session
where refresh_expires_at < now()
;
