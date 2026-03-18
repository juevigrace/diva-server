-- name: Count :one
SELECT COUNT(*) FROM diva_user;

-- name: CreateUser :exec
INSERT INTO diva_user (id, email, username, password_hash, birth_date, phone_number, alias, avatar, bio, user_verified, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: UpdateUser :exec
UPDATE diva_user 
SET alias = $1, avatar = $2, bio = $3, updated_at = NOW()
WHERE id = $4 AND deleted_at IS NULL;

-- name: UpdatePhoneNumber :exec
UPDATE diva_user
SET phone_number = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdateUsername :exec
UPDATE diva_user
SET username = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdateVerified :exec
UPDATE diva_user
SET user_verified = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdatePassword :exec
UPDATE diva_user
SET password_hash = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: UpdateEmail :exec
UPDATE diva_user
SET email = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE diva_user 
SET deleted_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT
  u.id,
  u.email,
  u.username,
  u.password_hash AS password_hash,
  u.alias,
  u.avatar,
  u.bio,
  u.user_verified AS user_verified,
  u.role,
  u.created_at AS created_at,
  u.updated_at AS updated_at,
  u.deleted_at AS deleted_at
FROM diva_user u
WHERE u.deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: GetUserByID :one
SELECT
  u.id,
  u.email,
  u.username,
  u.password_hash AS password_hash,
  u.alias,
  u.avatar,
  u.bio,
  u.user_verified AS user_verified,
  u.role,
  u.created_at AS created_at,
  u.updated_at AS updated_at,
  u.deleted_at AS deleted_at
FROM diva_user u
WHERE u.id = $1 AND u.deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT
  id,
  email,
  username,
  password_hash AS password_hash,
  alias,
  avatar,
  bio,
  user_verified AS user_verified,
  role,
  created_at AS created_at,
  updated_at AS updated_at,
  deleted_at AS deleted_at
FROM diva_user
WHERE username = $1 OR email = $2 AND deleted_at IS NULL;
