-- name: GetUserActionByID :one
select
    ua.id as id,
    ua.name,
    ua.user_id as userId
from diva_user_action ua
where ua.id = $1
;

-- name: GetUserActionByUserAndName :one
select
    ua.id as id,
    ua.name,
    ua.user_id as userId
from diva_user_action ua
where ua.user_id = $1 and ua.name = $2
;

-- name: CreateUserAction :exec
insert into diva_user_action (
    id,
    name,
    user_id
) values (
    $1,
    $2,
    $3
) on conflict (user_id, name) do nothing;

-- name: DeleteUserAction :exec
delete from diva_user_action
where id = $1;
