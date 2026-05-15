-- name: GetUserPermissions :many
select
    up.permission_id as permissionId,
    up.user_id as userId,
    up.granted_by as grantedBy,
    up.granted,
    up.granted_at as grantedAt,
    up.expires_at as expiresAt,
    up.updated_at as updatedAt,
    p.name,
    p.description,
    p.action,
    p.role_level as roleLevel
from diva_user_permissions up
join diva_permissions p on up.permission_id = p.id
where up.user_id = $1 and up.granted = true
;

-- name: GetUserPermission :one
select
    up.permission_id as permissionId,
    up.user_id as userId,
    up.granted_by as grantedBy,
    up.granted,
    up.granted_at as grantedAt,
    up.expires_at as expiresAt,
    up.updated_at as updatedAt
from diva_user_permissions up
where up.permission_id = $1 and up.user_id = $2
;

-- name: GrantPermission :exec
insert into diva_user_permissions (
    permission_id,
    user_id,
    granted_by,
    granted,
    granted_at,
    expires_at,
    updated_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) on conflict (permission_id, user_id) do update
set
    granted = $4,
    granted_by = $3,
    granted_at = $5,
    expires_at = $6,
    updated_at = now()
;

-- name: RevokePermission :exec
update diva_user_permissions
set
    granted = false,
    updated_at = now()
where permission_id = $1 and user_id = $2;

-- name: UpdatePermissionGrant :exec
update diva_user_permissions
set
    granted = $1,
    granted_by = $2,
    granted_at = now(),
    expires_at = $3,
    updated_at = now()
where permission_id = $4 and user_id = $5;

-- name: DeleteUserPermission :exec
delete from diva_user_permissions
where permission_id = $1 and user_id = $2;

-- name: DeleteAllUserPermissions :exec
delete from diva_user_permissions
where user_id = $1;