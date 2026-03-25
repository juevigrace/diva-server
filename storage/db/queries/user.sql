-- name: Count :one
select count(*)
from diva_user
;

-- name: CreateUser :exec
INSERT INTO diva_user (id, email, username, password_hash, alias)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateProfile :exec
UPDATE diva_user 
SET alias = $1, avatar = $2, bio = $3, updated_at = NOW()
WHERE id = $4 AND deleted_at IS NULL;

-- name: UpdatePhoneNumber :exec
UPDATE diva_user
SET phone_number = $1,
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdateUsername :exec
UPDATE diva_user
SET username = $1,
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdateVerified :exec
UPDATE diva_user
SET user_verified = $1,
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdatePassword :exec
UPDATE diva_user
SET password_hash = $1,
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdateEmail :exec
UPDATE diva_user
SET email = $1,
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE diva_user 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAllUsers :many
select
    u.id,
    u.email,
    u.username,
    u.password_hash,
    u.phone_number,
    u.birth_date,
    u.alias,
    u.avatar,
    u.bio,
    u.user_verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.deleted_at is null
limit $1
offset $2
;

-- name: GetUserByID :one
select
    u.id,
    u.email,
    u.username,
    u.password_hash,
    u.phone_number,
    u.birth_date,
    u.alias,
    u.avatar,
    u.bio,
    u.user_verified,
    u.role,
    u.created_at,
    u.updated_at,
    u.deleted_at
from diva_user u
where u.id = $1 and u.deleted_at is null
;

-- name: GetUserByUsername :one
select
    id,
    email,
    username,
    password_hash,
    phone_number,
    birth_date,
    alias,
    avatar,
    bio,
    user_verified,
    role,
    created_at,
    updated_at,
    deleted_at
from diva_user
where username = $1 and deleted_at is null
;

-- name: GetUserByEmail :one
select
    id,
    email,
    username,
    password_hash,
    phone_number,
    birth_date,
    alias,
    avatar,
    bio,
    user_verified,
    role,
    created_at,
    updated_at,
    deleted_at
from diva_user
where email = $1 and deleted_at is null
;

-- name: GetUserByUsernameOrEmail :one
select
    id,
    email,
    username,
    password_hash,
    phone_number,
    birth_date,
    alias,
    avatar,
    bio,
    user_verified,
    role,
    created_at,
    updated_at,
    deleted_at
from diva_user
where (username = $1 or email = $2) and deleted_at is null
;
