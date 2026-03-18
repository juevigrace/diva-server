-- name: CountPermissions :one
SELECT COUNT(*)
FROM diva_permissions
WHERE deleted_at IS NULL;

-- name: CreatePermission :exec
INSERT INTO diva_permissions (id, name, description, role_level)
VALUES ($1, $2, $3, $4);

-- name: UpdatePermission :exec
UPDATE diva_permissions
SET name = $1, description = $2, updated_at = NOW()
WHERE id = $3 AND deleted_at IS NULL;

-- name: DeletePermission :exec
UPDATE diva_permissions
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListPermissions :many
SELECT
  id,
  name,
  description,
  role_level AS role_level,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_permissions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPermissionByID :one
SELECT
  id,
  name,
  description,
  role_level AS role_level,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_permissions
WHERE id = $1;
