-- name: GetSessionByID :one
select
    s.id as id,
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

-- name: ListSessionsByUser :many
select
    s.id as id,
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
update diva_session set
    access_token = $1,
    refresh_token = $2,
    ip_address = $3,
    expires_at = $4,
    updated_at = now()
where id = $5;

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
where expires_at < now()
;
