-- name: CreateUserAction :exec
INSERT INTO diva_user_pending_actions (id, user_id, action_name)
VALUES ($1, $2, $3);

-- name: GetUserActions :many
select *
from diva_user_pending_actions
where user_id = $1
;

-- name: GetUserAction :one
select *
from diva_user_pending_actions
where action_name = $1 and user_id = $2
;

-- name: DeleteUserAction :exec
delete from diva_user_pending_actions
where user_id = $1 and action_name = $2
;
