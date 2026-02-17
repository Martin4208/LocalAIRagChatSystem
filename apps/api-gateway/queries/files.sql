-- name: GetFileByHash :one
SELECT * FROM files 
WHERE sha256_hash = $1 
LIMIT 1;

-- name: CreateFile :one
INSERT INTO files (
    sha256_hash, 
    mime_type, 
    size_bytes, 
    original_filename, 
    minio_bucket, 
    minio_key,
    created_at,
    last_accessed_at
) VALUES (
    $1, $2, $3, $4, $5, $6, now(), now()
)
RETURNING *;

-- name: GetFileById :one
SELECT * FROM files 
WHERE id = $1;

-- name: UpdateFileLastAccessed :exec
UPDATE files 
SET last_accessed_at = now() 
WHERE id = $1;