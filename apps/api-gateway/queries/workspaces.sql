-- name: CreateWorkspace :one
INSERT INTO workspaces (
    name, description, settings
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetWorkspace :one
SELECT * FROM workspaces WHERE id =$1 AND deleted_at IS NULL;

-- name: ListWorkspaces :many
SELECT * FROM workspaces WHERE deleted_at IS NULL;

-- name: UpdateWorkspace :exec
UPDATE workspaces 
SET 
    name = $1, 
    description = $2,
    updated_at = now() 
WHERE id = $3 AND deleted_at IS NULl;

-- name: DeleteWorkspace :exec
UPDATE workspaces
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;
