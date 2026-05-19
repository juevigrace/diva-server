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
where up.id = $1
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
where up.user_id = $1
;

-- name: GetPreferencesByDevice :one
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
where up.device = $1
;

-- name: CreateUserPreferences :exec
insert into diva_user_preferences (
    id,
    user_id,
    device,
    theme,
    onboarding_completed,
    language,
    created_at,
    updated_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: UpdateUserPreferences :exec
update diva_user_preferences set
    theme = $1,
    language = $2,
    last_sync_at = now(),
    updated_at = $3
where id = $4;

-- name: DeletePreferences :exec
delete from diva_user_preferences
where id = $1
;

-- name: DeletePreferencesByUser :exec
delete from diva_user_preferences
where user_id = $1
;
