-- name: GetUserStateByUserID :one
select us.user_id, us.verified, us.status, us.last_active_at, us.updated_at
from diva_user_state us
where us.user_id = $1
;

-- name: CreateUserState :exec
insert into diva_user_state (
    user_id,
    verified,
    status
) values (
    $1,
    $2,
    $3
);

-- name: UpdateUserVerified :exec
update diva_user_state set
    verified = $1,
    updated_at = now()
where user_id = $2;

-- name: UpdateUserStatus :exec
update diva_user_state set
    status = $1,
    updated_at = now()
where user_id = $2;

-- name: UpdateLastActiveAt :exec
update diva_user_state set
    last_active_at = now(),
    updated_at = now()
where user_id = $1;
