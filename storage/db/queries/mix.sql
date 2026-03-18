-- name: CreateMix :exec
INSERT INTO diva_mix_metadata (collection_id, algorithm_type, time_window_hours, content_weight, freshness_weight, min_engagement_score, excluded_tags, auto_refresh, refresh_interval_seconds)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: UpdateMix :exec
UPDATE diva_mix_metadata
SET algorithm_type = $1, time_window_hours = $2, content_weight = $3, freshness_weight = $4, min_engagement_score = $5, excluded_tags = $6, auto_refresh = $7, refresh_interval_seconds = $8
WHERE collection_id = $9;

-- name: DeleteMix :exec
DELETE FROM diva_mix_metadata
WHERE collection_id = $1;
