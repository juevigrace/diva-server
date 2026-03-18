-- Post Interactions

-- name: CreatePostInteraction :exec
INSERT INTO diva_post_interaction (id, post_id, user_id, reaction_type)
VALUES ($1, $2, $3, $4);

-- name: DeletePostInteraction :exec
DELETE FROM diva_post_interaction
WHERE id = $1;

-- Collection Interactions

-- name: CreateCollectionInteraction :exec
INSERT INTO diva_collection_interaction (id, collection_id, user_id, reaction_type)
VALUES ($1, $2, $3, $4);

-- name: DeleteCollectionInteraction :exec
DELETE FROM diva_collection_interaction
WHERE id = $1;

-- Message Interactions

-- name: CreateMessageInteraction :exec
INSERT INTO diva_message_interaction (id, message_id, user_id, reaction_type)
VALUES ($1, $2, $3, $4);

-- name: DeleteMessageInteraction :exec
DELETE FROM diva_message_interaction
WHERE id = $1;

-- Post Shares

-- name: CreatePostShare :exec
INSERT INTO diva_post_share (interaction_id, share_type, caption)
VALUES ($1, $2, $3);

-- name: DeletePostShare :exec
DELETE FROM diva_post_share
WHERE interaction_id = $1;

-- Post Comments

-- name: CreatePostComment :exec
INSERT INTO diva_post_comment (interaction_id, reply_to, content)
VALUES ($1, $2, $3);

-- name: UpdatePostComment :exec
UPDATE diva_post_comment
SET content = $1, edited_at = NOW()
WHERE interaction_id = $2;

-- name: DeletePostComment :exec
DELETE FROM diva_post_comment
WHERE interaction_id = $1;

-- Collection Shares

-- name: CreateCollectionShare :exec
INSERT INTO diva_collection_share (interaction_id, share_type, caption)
VALUES ($1, $2, $3);

-- name: DeleteCollectionShare :exec
DELETE FROM diva_collection_share
WHERE interaction_id = $1;

-- Collection Comments

-- name: CreateCollectionComment :exec
INSERT INTO diva_collection_comment (interaction_id, reply_to, content)
VALUES ($1, $2, $3);

-- name: UpdateCollectionComment :exec
UPDATE diva_collection_comment
SET content = $1, edited_at = NOW()
WHERE interaction_id = $2;

-- name: DeleteCollectionComment :exec
DELETE FROM diva_collection_comment
WHERE interaction_id = $1;
