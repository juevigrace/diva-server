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
where up.user_id = $1 and up.granted = true
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
where up.permission_id = $1 and up.user_id = $2
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
where up.user_id = $1 and dp.name = $2
;

-- name: CreateUserPermission :exec
insert into diva_user_permissions (
    permission_id,
    user_id,
    granted_by,
    granted,
    expires_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5
)
;

-- name: UpdateUserPermission :exec
update diva_user_permissions
set
    granted = $1,
    expires_at = $2,
    updated_at = now()
where permission_id = $3 and user_id = $4;

-- name: DeleteUserPermission :exec
delete from diva_user_permissions
where permission_id = $1 and user_id = $2
;
