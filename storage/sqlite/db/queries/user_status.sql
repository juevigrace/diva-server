-- name: GetUserStateByUserID :one
select us.user_id, us.verified, us.status, us.last_active_at, us.updated_at
from diva_user_state us
where us.user_id = ?
;

-- name: CreateUserState :exec
insert into diva_user_state (
    user_id,
    verified,
    status
) values (
    ?,
    ?,
    ?
);

-- name: UpdateUserVerified :exec
update diva_user_state set
    verified = ?,
    updated_at = CURRENT_TIMESTAMP
where user_id = ?;

-- name: UpdateUserStatus :exec
update diva_user_state set
    status = ?,
    updated_at = CURRENT_TIMESTAMP
where user_id = ?;

-- name: UpdateLastActiveAt :exec
update diva_user_state set
    last_active_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
where user_id = ?;
