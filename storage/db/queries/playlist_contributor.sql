-- name: CreatePlaylistContributor :exec
INSERT INTO diva_playlist_contributor (collection_id, contributor_id)
VALUES ($1, $2);

-- name: UpdatePlaylistContributor :exec
UPDATE diva_playlist_contributor
SET collection_id = $1, contributor_id = $2
WHERE collection_id = $3 AND contributor_id = $4;

-- name: DeletePlaylistContributor :exec
DELETE FROM diva_playlist_contributor
WHERE collection_id = $1 AND contributor_id = $2;
