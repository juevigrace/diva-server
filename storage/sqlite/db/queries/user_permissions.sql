-- name: GetUserPermissions :many
select
    up.permission_id,
    up.user_id,
    up.granted_by,
    up.granted,
    up.granted_at,
    up.expires_at,
    up.updated_at
from diva_user_permissions up
where up.user_id = ? and up.granted = true
;

-- name: GetUserPermission :one
select
    up.permission_id,
    up.user_id,
    up.granted_by,
    up.granted,
    up.granted_at,
    up.expires_at,
    up.updated_at
from diva_user_permissions up
where up.permission_id = ? and up.user_id = ?
;

-- name: GetUserPermissionByName :one
select
    up.permission_id,
    up.user_id,
    up.granted_by,
    up.granted,
    up.granted_at,
    up.expires_at,
    up.updated_at
from diva_user_permissions up
left join diva_permissions dp on up.permission_id = dp.id
where up.user_id = ? and dp.name = ?
;

-- name: CreateUserPermission :exec
insert into diva_user_permissions (
    permission_id,
    user_id,
    granted_by,
    granted,
    expires_at
) values (
    ?,
    ?,
    ?,
    ?,
    ?
)
;

-- name: UpdateUserPermission :exec
update diva_user_permissions
set
    granted = ?,
    expires_at = ?,
    updated_at = CURRENT_TIMESTAMP
where permission_id = ? and user_id = ?;

-- name: DeleteUserPermission :exec
delete from diva_user_permissions
where permission_id = ? and user_id = ?
;

-- name: DeleteExpiredUserPermissions :exec
delete from diva_user_permissions
where expires_at is not null and expires_at < CURRENT_TIMESTAMP
;
