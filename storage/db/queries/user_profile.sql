-- name: GetUserProfileByUserID :one
select
    up.user_id as userId,
    up.first_name as firstName,
    up.last_name as lastName,
    up.birth_date as birthDate,
    up.phone_number as phoneNumber,
    up.alias,
    up.bio
from diva_user_profile up
where up.user_id = $1
;

-- name: CreateUserProfile :exec
insert into diva_user_profile (
    user_id,
    first_name,
    last_name,
    birth_date,
    phone_number,
    alias,
    bio
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);

-- name: UpdateUserProfile :exec
update diva_user_profile
set
    first_name = $1,
    last_name = $2,
    birth_date = $3,
    alias = $4,
    bio = $5
where user_id = $6;

-- name: UpdatePhoneNumber :exec
update diva_user_profile
set
    phone_number = $1
where user_id = $2;

-- name: DeleteUserProfile :exec
delete from diva_user_profile
where user_id = $1;