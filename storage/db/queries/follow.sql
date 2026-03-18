-- name: CreateFollow :exec
INSERT INTO diva_follow (user_id, followed)
VALUES ($1, $2);

-- name: DeleteFollow :exec
DELETE FROM diva_follow
WHERE user_id = $1 AND followed = $2;
