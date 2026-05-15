-- name: GetActionVerification :one
select 
    v.token,
    v.expires_at as expiresat
from diva_user_action_verification v
where v.action_id = $1 and v.token = $2
;

-- name: CreateVerification :exec
insert into diva_user_action_verification (
    action_id,
    token,
    expires_at
) values (
    $1,
    $2,
    $3
);

-- name: DeleteVerification :exec
delete from diva_user_action_verification
where action_id = $1 and token = $2
;

-- name: DeleteVerificationByActionID :exec
delete from diva_user_action_verification
where action_id = $1
;
