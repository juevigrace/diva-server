-- name: GetUserProfileByUserID :one
select
    up.user_id, up.first_name, up.last_name, up.birth_date, up.alias, up.bio, up.avatar
from diva_user_profile up
where up.user_id = $1
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
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateUserProfile :exec
update diva_user_profile set
    first_name = $1,
    last_name = $2,
    birth_date = $3,
    alias = $4,
    bio = $5
where user_id = $6;

-- name: UpdateUserProfileAvatar :exec
update diva_user_profile set
    avatar = $1
where user_id = $2;

-- name: DeleteUserProfile :exec
delete from diva_user_profile
where user_id = $1
;
