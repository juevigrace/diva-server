-- name: GetUserVerification :one
select v.action_id, v.token, v.expires_at, v.used_at, v.verified
from diva_action_verification v
where v.action_id = $1
;

-- name: CreateUserVerification :exec
insert into diva_action_verification (
    action_id,
    token,
    expires_at
) values (
    $1,
    $2,
    $3
);

-- name: UpdateUserVerification :exec
update diva_action_verification set
    verified = $1,
    used_at = now()
where action_id = $2;

-- name: DeleteUserVerification :exec
delete from diva_action_verification
where action_id = $1
;
