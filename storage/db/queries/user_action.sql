-- name: ListActionsByUser :many
select ua.id, ua.name, ua.user_id
from diva_action ua
where ua.user_id = $1
;

-- name: GetUserActionByID :one
select ua.id, ua.name, ua.user_id
from diva_action ua
where ua.id = $1
;

-- name: GetUserActionByUserAndName :one
select ua.id, ua.name, ua.user_id
from diva_action ua
where ua.user_id = $1 and ua.name = $2
;

-- name: CreateUserAction :exec
insert into diva_action (
    id,
    name,
    user_id
) values (
    $1,
    $2,
    $3
) on conflict (user_id, name) do nothing;

-- name: DeleteUserAction :exec
delete from diva_action
where id = $1
;

-- name: DeleteUserActionByUser :exec
delete from diva_action
where user_id = $1
;
