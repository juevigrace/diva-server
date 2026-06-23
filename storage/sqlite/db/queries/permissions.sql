-- name: GetPermissionByID :one
select
    p.id as id,
    p.name,
    p.description,
    p.action,
    p.role_level,
    p.created_at,
    p.updated_at,
    p.deleted_at
from diva_permissions p
where p.id = ?
;

-- name: GetPermissionByName :one
select
    p.id as id,
    p.name,
    p.description,
    p.action,
    p.role_level,
    p.created_at,
    p.updated_at,
    p.deleted_at
from diva_permissions p
where p.name = ?
;

-- name: ListPermissions :many
select
    p.id as id,
    p.name,
    p.description,
    p.action,
    p.role_level,
    p.created_at,
    p.updated_at,
    p.deleted_at
from diva_permissions p
order by p.name
limit ?
offset ?
;

-- name: CountPermissions :one
select count(*)
from diva_permissions
;

-- name: CreatePermission :exec
insert into diva_permissions (
    id,
    name,
    description,
    action,
    role_level
) values (
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdatePermission :exec
update diva_permissions set
    name = ?,
    description = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdatePermissionAction :exec
update diva_permissions set
    action = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdatePermissionRoleLevel :exec
update diva_permissions set
    role_level = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: DeletePermission :exec
delete from diva_permissions
where id = ?
;

-- name: SoftDeletePermission :exec
update diva_permissions
set
    deleted_at = CURRENT_TIMESTAMP
where id = ?;

-- name: RestorePermission :exec
update diva_permissions set
    deleted_at = null
where id = ?;
