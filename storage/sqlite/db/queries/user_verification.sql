-- name: GetUserVerification :one
select v.action_id, v.token, v.expires_at, v.used_at, v.verified
from diva_action_verification v
where v.action_id = ?
;

-- name: CreateUserVerification :exec
insert into diva_action_verification (
    action_id,
    token,
    expires_at
) values (
    ?,
    ?,
    ?
);

-- name: UpdateUserVerification :exec
update diva_action_verification set
    verified = ?,
    used_at = CURRENT_TIMESTAMP
where action_id = ?;

-- name: DeleteUserVerification :exec
delete from diva_action_verification
where action_id = ?
;
