-- name: GetEmailVerificationTokenByToken :one
SELECT
    user_id AS user_id,
    token,
    expires_at AS expires_at,
    created_at AS created_at,
    used_at AS used_at
FROM diva_email_verification_tokens
WHERE token = $1;

-- name: GetEmailVerificationTokenByUserID :one
SELECT
    user_id AS user_id,
    token,
    expires_at AS expires_at,
    created_at AS created_at,
    used_at AS used_at
FROM diva_email_verification_tokens
WHERE user_id = $1;

-- name: CreateEmailVerificationToken :exec
INSERT INTO diva_email_verification_tokens(
    user_id,
    token,
    expires_at
) VALUES ($1, $2, $3);

-- name: MarkEmailVerificationTokenAsUsed :exec
UPDATE diva_email_verification_tokens
SET used_at = NOW()
WHERE user_id = $1;

-- name: DeleteEmailVerificationTokenByToken :exec
DELETE FROM diva_email_verification_tokens
WHERE token = $1;

-- name: DeleteEmailVerificationToken :exec
DELETE FROM diva_email_verification_tokens
WHERE user_id = $1;
