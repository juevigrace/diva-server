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
where p.id = $1
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
where p.action = $1
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
limit $1
offset $2
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
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: UpdatePermission :exec
update diva_permissions set
    name = $1,
    description = $2,
    updated_at = now()
where id = $3;

-- name: UpdatePermissionAction :exec
update diva_permissions set
    action = $1,
    updated_at = now()
where id = $2;

-- name: UpdatePermissionRoleLevel :exec
update diva_permissions set
    role_level = $1,
    updated_at = now()
where id = $2;

-- name: DeletePermission :exec
delete from diva_permissions
where id = $1
;

-- name: SoftDeletePermission :exec
update diva_permissions
set
    deleted_at = now()
where id = $1;

-- name: RestorePermission :exec
update diva_permissions set
    deleted_at = null
where id = $1;
