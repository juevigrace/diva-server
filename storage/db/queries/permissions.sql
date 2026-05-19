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
where p.name = $1
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
select count(*) from diva_permissions;

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
    action = $3,
    role_level = $4,
    updated_at = now()
where id = $5;

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
