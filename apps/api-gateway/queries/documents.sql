-- name: CreateDocument :one
INSERT INTO documents (
    workspace_id, 
    directory_id, 
    file_id, 
    name, 
    tags, 
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, 'uploaded', now(), now()
)
RETURNING *;

-- name: GetDocument :one
SELECT 
    d.id,
    d.workspace_id,
    d.directory_id,
    d.file_id,
    d.name,
    d.tags,
    d.status,
    d.created_at,
    d.updated_at,
    f.size_bytes,
    f.mime_type,
    f.sha256_hash,
    f.minio_bucket,
    f.minio_key
FROM documents d
INNER JOIN files f ON d.file_id = f.id
WHERE d.id = $1 
  AND d.workspace_id = $2 
  AND d.deleted_at IS NULL;

-- name: ListDocuments :many
SELECT 
    d.id,
    d.workspace_id,
    d.directory_id,
    d.file_id,
    d.name,
    d.tags,
    d.status,
    d.created_at,
    d.updated_at,
    f.size_bytes,
    f.mime_type
FROM documents d
INNER JOIN files f ON d.file_id = f.id
WHERE d.workspace_id = $1 
  AND d.deleted_at IS NULL
  AND (sqlc.narg('directory_id')::uuid IS NULL OR d.directory_id = sqlc.narg('directory_id'))
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountDocuments :one
SELECT COUNT(*) 
FROM documents 
WHERE workspace_id = $1 
  AND deleted_at IS NULL
  AND (sqlc.narg('directory_id')::uuid IS NULL OR directory_id = sqlc.narg('directory_id'));

-- name: DeleteDocument :exec
UPDATE documents 
SET 
    deleted_at = now(),
    updated_at = now()
WHERE id = $1 
  AND workspace_id = $2;