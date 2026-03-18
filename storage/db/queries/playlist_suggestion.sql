-- name: CreatePlaylistSuggestion :exec
INSERT INTO diva_playlist_suggestions (collection_id, suggester_id, media_id, status)
VALUES ($1, $2, $3, $4);

-- name: UpdatePlaylistSuggestion :exec
UPDATE diva_playlist_suggestions
SET status = $1
WHERE collection_id = $2 AND suggester_id = $3 AND media_id = $4;
