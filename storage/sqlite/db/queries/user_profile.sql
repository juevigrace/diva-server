-- name: GetUserProfileByUserID :one
select
    up.user_id,
    up.first_name,
    up.last_name,
    up.birth_date,
    up.alias,
    up.bio,
    up.avatar,
    up.updated_at
from diva_user_profile up
where up.user_id = ?
;

-- name: CreateUserProfile :exec
insert into diva_user_profile (
    user_id,
    first_name,
    last_name,
    birth_date,
    alias,
    bio
) values (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdateUserProfile :exec
update diva_user_profile set
    first_name = ?,
    last_name = ?,
    birth_date = ?,
    alias = ?,
    bio = ?,
    updated_at = CURRENT_TIMESTAMP
where user_id = ?;

-- name: UpdateUserProfileAvatar :exec
update diva_user_profile set
    avatar = ?,
    updated_at = CURRENT_TIMESTAMP
where user_id = ?;
