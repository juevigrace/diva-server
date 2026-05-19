-- name: GetVerification :one
select v.action_id, v.token, v.expires_at
from diva_action_verification v
where v.action_id = $1
;

-- name: CreateVerification :exec
insert into diva_action_verification (
    action_id,
    token,
    expires_at
) values (
    $1,
    $2,
    $3
);

-- name: DeleteVerification :exec
delete from diva_action_verification
where action_id = $1
;
