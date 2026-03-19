-- name: Create :exec
INSERT INTO diva_email_verification_tokens (user_id, token, expires_at, created_at)
VALUES ($1, $2, $3, $4);

-- name: GetByToken :one
select user_id, token, expires_at, created_at
from diva_email_verification_tokens
where token = $1
;

-- name: DeleteByToken :exec
delete from diva_email_verification_tokens
where token = $1
;
