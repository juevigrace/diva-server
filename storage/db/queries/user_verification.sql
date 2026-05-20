-- name: GetVerification :one
select v.action_id, v.token, v.expires_at, v.used_at, v.verified
from diva_action_verification v
where v.action_id = $1
;

-- name: CreateVerification :exec
insert into diva_action_verification (
    action_id,
    token,
    verified,
    expires_at,
    used_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: UpdateVerification :exec
update diva_action_verification set
    verified = $1,
    used_at = now()
where action_id = $2;

-- name: DeleteVerification :exec
delete from diva_action_verification
where action_id = $1
;
