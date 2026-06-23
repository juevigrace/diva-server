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
    s.expires_at,
    s.created_at,
    s.updated_at
from diva_session s
where s.id = ?
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
    s.expires_at,
    s.created_at,
    s.updated_at
from diva_session s
where s.user_id = ?
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
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdateSession :exec
update diva_session set
    access_token = ?,
    refresh_token = ?,
    ip_address = ?,
    expires_at = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdateSessionStatus :exec
update diva_session set
    status = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: DeleteSession :exec
delete from diva_session
where id = ?
;

-- name: DeleteSessionsByUser :exec
delete from diva_session
where user_id = ?
;

-- name: DeleteExpiredSessions :exec
delete from diva_session
where expires_at < CURRENT_TIMESTAMP
;
