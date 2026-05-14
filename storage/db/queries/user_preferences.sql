-- name: CreateUserPreferences :exec
INSERT INTO diva_user_preferences (id, user_id, device, theme, onboarding_completed, language, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateUserPreferences :exec
UPDATE diva_user_preferences 
SET theme = $1, language = $2, last_sync_at = NOW(), updated_at = $3
WHERE id = $4;

-- name: GetPreferencesByID :one
select
    up.id,
    up.user_id,
    up.device,
    up.theme,
    up.onboarding_completed,
    up.language,
    up.last_sync_at,
    up.created_at,
    up.updated_at
from diva_user_preferences as up
where up.id = $1
;

-- name: GetPreferencesByUser :many
select
    up.id,
    up.user_id,
    up.device,
    up.theme,
    up.onboarding_completed,
    up.language,
    up.last_sync_at,
    up.created_at,
    up.updated_at
from diva_user_preferences as up
where up.user_id = $1
;

-- name: GetPreferencesByDevice :one
select
    up.id,
    up.user_id,
    up.device,
    up.theme,
    up.onboarding_completed,
    up.language,
    up.last_sync_at,
    up.created_at,
    up.updated_at
from diva_user_preferences as up
where up.device = $1
;
