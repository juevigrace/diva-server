-- name: CreateVerification :exec
INSERT INTO diva_email_verification_tokens (user_id, action_id, token, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetVerificationByToken :one
select
    ev.user_id,
    ev.token,
    ev.expires_at,
    ev.created_at,
    up.id as action_id,
    up.action_name
from diva_email_verification_tokens as ev
left join diva_user_pending_actions as up on up.id = ev.action_id
where token = $1
;

-- name: DeleteByToken :exec
delete from diva_email_verification_tokens
where token = $1
;
