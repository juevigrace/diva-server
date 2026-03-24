-- name: CreateUserPendingAction :exec
INSERT INTO diva_user_pending_actions (user_id, action_name)
VALUES ($1, $2);

-- name: GetUserPendingActions :many
SELECT user_id, action_name
FROM diva_user_pending_actions
WHERE user_id = $1;

-- name: DeleteUserPendingAction :exec
DELETE FROM diva_user_pending_actions
WHERE user_id = $1 AND action_name = $2;
