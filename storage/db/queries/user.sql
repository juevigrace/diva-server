-- name: GetUserByID :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.id = $1
;

-- name: GetUserByUsernameOrEmail :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.email = $1 or u.username = $1
;

-- name: GetUserByEmail :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.email = $1
;

-- name: GetUserByUsername :one
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.username = $1
;

-- name: ListUsers :many
select
    u.id as id,
    u.username,
    u.email,
    u.phone_number,
    u.password_hash,
    u.verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
order by u.created_at desc
limit $1
offset $2
;

-- name: CreateUser :exec
insert into diva_user (
    id,
    username,
    email,
    password_hash,
    verified,
    role
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateUsername :exec
update diva_user set
    username = $1,
    updated_at = now()
where id = $2;

-- name: UpdateEmail :exec
update diva_user set
    email = $1,
    updated_at = now()
where id = $2;

-- name: UpdatePassword :exec
update diva_user set
    password_hash = $1,
    updated_at = now()
where id = $2;

-- name: UpdateVerified :exec
update diva_user set
    verified = $1,
    updated_at = now()
where id = $2;

-- name: UpdateRole :exec
update diva_user set
    role = $1,
    updated_at = now()
where id = $2;

-- name: UpdatePhoneNumber :exec
update diva_user set
    phone_number = $1,
    updated_at = now()
where id = $2;

-- name: DeleteUser :exec
delete from diva_user
where id = $1
;

-- name: CountUsers :one
select count(*)
from diva_user
;

-- name: SoftDeleteUser :exec
update diva_user set
    deleted_at = now()
where id = $1;

-- name: RestoreUser :exec
update diva_user set
    deleted_at = null
where id = $1;
