-- name: ListActionsByUser :many
select ua.id, ua.name, ua.user_id
from diva_action ua
where ua.user_id = ?
;

-- name: GetUserActionByID :one
select ua.id, ua.name, ua.user_id
from diva_action ua
where ua.id = ?
;

-- name: GetUserActionByUserAndName :one
select ua.id, ua.name, ua.user_id
from diva_action ua
where ua.user_id = ? and ua.name = ?
;

-- name: CreateUserAction :exec
insert into diva_action (
    id,
    name,
    user_id
) values (
    ?,
    ?,
    ?
);

-- name: DeleteUserAction :exec
delete from diva_action
where id = ?
;

-- name: DeleteUserActionByUser :exec
delete from diva_action
where user_id = ?
;

-- name: DeleteExpiredActions :exec
delete from diva_action
where id in (
    select action_id from diva_action_verification
    where expires_at < CURRENT_TIMESTAMP and verified = false
)
;
