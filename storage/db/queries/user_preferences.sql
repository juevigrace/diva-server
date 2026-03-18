-- name: CreateUserPreferences :exec
INSERT INTO diva_user_preferences (id, user_id, theme, onboarding_completed, language, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateUserPreferences :exec
UPDATE diva_user_preferences 
SET theme = $1, onboarding_completed = $2, language = $3, last_sync_at = NOW(), updated_at = $4
WHERE id = $5;
