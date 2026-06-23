-- name: GetPreferencesByID :one
select
    up.id as id,
    up.user_id,
    up.device,
    up.theme,
    up.onboarding_completed,
    up.language,
    up.last_sync_at,
    up.created_at,
    up.updated_at
from diva_user_preferences up
where up.id = ?
;

-- name: GetPreferencesByUser :many
select
    up.id as id,
    up.user_id,
    up.device,
    up.theme,
    up.onboarding_completed,
    up.language,
    up.last_sync_at,
    up.created_at,
    up.updated_at
from diva_user_preferences up
where up.user_id = ?
;

-- name: CreateUserPreferences :exec
insert into diva_user_preferences (
    id,
    user_id,
    device,
    theme,
    onboarding_completed,
    language
) values (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdateUserPreferences :exec
update diva_user_preferences set
    theme = ?,
    language = ?,
    last_sync_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
where id = ?;
