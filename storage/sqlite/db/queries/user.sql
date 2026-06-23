-- name: GetUserByID :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.id = ?
;

-- name: GetUserByUsernameOrEmail :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.email = ? or u.username = ?
;

-- name: GetUserByEmail :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.email = ?
;

-- name: GetUserByUsername :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.username = ?
;

-- name: ListUsers :many
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
order by u.created_at desc
limit ?
offset ?
;

-- name: CreateUser :exec
insert into diva_user (
    id,
    username,
    email,
    password_hash,
    role
) values (
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdateUsername :exec
update diva_user set
    username = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdateEmail :exec
update diva_user set
    email = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdatePassword :exec
update diva_user set
    password_hash = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdatePhoneNumber :exec
update diva_user set
    phone_number = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: UpdateRole :exec
update diva_user set
    role = ?,
    updated_at = CURRENT_TIMESTAMP
where id = ?;

-- name: DeleteUser :exec
delete from diva_user
where id = ?
;

-- name: CountUsers :one
select count(*)
from diva_user
;

-- name: SoftDeleteUser :exec
update diva_user set
    deleted_at = CURRENT_TIMESTAMP
where id = ?;

-- name: RestoreUser :exec
update diva_user set
    deleted_at = null,
    updated_at = CURRENT_TIMESTAMP
where id = ?;
