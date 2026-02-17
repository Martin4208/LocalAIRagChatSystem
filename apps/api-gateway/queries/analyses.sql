-- ========================================
-- Analysis Operations
-- ========================================

-- name: CreateAnalysis :one
INSERT INTO analyses (
    workspace_id,
    title,
    description,
    analysis_type,
    status,
    config,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, 'pending', $5, now(), now()
)
RETURNING *;

-- name: ListAnalyses :many
SELECT 
    id,
    workspace_id,
    title,
    description,
    analysis_type,
    status,
    started_at,
    completed_at,
    config,
    error_message,
    created_at,
    updated_at
FROM analyses
WHERE 
    workspace_id = $1
    AND deleted_at IS NULL
    AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountAnalyses :one
SELECT COUNT(*)
FROM analyses
WHERE
    workspace_id = $1
    AND deleted_at IS NULL
    AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'));

-- name: GetAnalysis :one
SELECT * FROM analyses
WHERE
    id = $1
    AND workspace_id = $2
    AND deleted_at IS NULL;

-- name: UpdateAnalysisStatus :exec
UPDATE analyses
SET
    status = $2,
    started_at = CASE
        WHEN $2 = 'processing' THEN now()
        ELSE started_at
    END,
    completed_at = CASE
        WHEN $2 IN ('completed', 'failed') THEN now()
        ELSE completed_at
    END,
    error_message = $3,
    updated_at = now()
WHERE id = $1;

-- name: DeleteAnalysis :exec
UPDATE analyses
SET
    deleted_at = now(),
    updated_at = now()
WHERE 
    id = $1
    AND workspace_id = $2
    AND deleted_at IS NULL;

-- ========================================
-- Analysis Results Operations
-- ========================================

-- name: CreateAnalysisResult :one
INSERT INTO analysis_results (
    analysis_id,
    result_type,
    content,
    image_url,
    minio_bucket,
    minio_key,
    metadata,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now()
)
RETURNING *;

-- name: ListAnalysisResults :many
SELECT
    id,
    analysis_id,
    result_type,
    content,
    image_url,
    minio_bucket,
    minio_key,
    metadata,
    created_at
FROM analysis_results
WHERE 
    analysis_id = $1
ORDER BY created_at ASC;

-- name: GetAnalysisResult :one
SELECT * FROM analysis_results
WHERE 
    id = $1
    AND analysis_id = $2;

-- name: DeleteAnalysisResults :exec
DELETE FROM analysis_results
WHERE analysis_id = $1;