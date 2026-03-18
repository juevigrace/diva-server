-- name: CreateUserPermission :exec
INSERT INTO diva_user_permissions (permission_id, user_id, granted_by, granted, granted_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUserPermission :exec
UPDATE diva_user_permissions
SET granted = $1, expires_at = $2, updated_at = NOW()
WHERE permission_id = $3 AND user_id = $4;

-- name: DeleteUserPermission :exec
DELETE FROM diva_user_permissions
WHERE permission_id = $1 AND user_id = $2;
