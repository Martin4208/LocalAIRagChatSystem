-- name: ListGraphsByWorkspace :many
SELECT 
    id,
    workspace_id,
    title,
    graph_type,
    created_at,
    updated_at
FROM graphs
WHERE workspace_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateGraph :one
INSERT INTO graphs (
    workspace_id,
    title,
    graph_type
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetGraphByID :one
SELECT 
    id,
    workspace_id,
    title,
    graph_type,
    layout_config,
    created_at,
    updated_at
FROM graphs
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateGraph :one
UPDATE graphs
SET 
    title = COALESCE(sqlc.narg('title')::text, title),
    layout_config = COALESCE(sqlc.narg('layout_config')::jsonb, layout_config),
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteGraph :exec
UPDATE graphs
SET deleted_at = now()
WHERE id = $1;

-- name: GetGraphNodesByGraphID :many
SELECT 
    id,
    graph_id,
    label,
    node_type,
    source_type,
    source_id,
    position,
    style,
    metadata,
    created_at
FROM graph_nodes
WHERE graph_id = $1
ORDER BY created_at ASC;

-- name: GetGraphEdgesByGraphID :many
SELECT 
    id,
    graph_id,
    from_node_id,
    to_node_id,
    edge_type,
    is_directed,
    weight,
    confidence,
    style,
    metadata,
    created_at,
    updated_at
FROM graph_edges
WHERE graph_id = $1
ORDER BY created_at ASC;

-- name: CreateGraphNode :one
INSERT INTO graph_nodes (
    graph_id,
    label,
    node_type,
    source_type,
    source_id,
    style,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, graph_id, label, node_type, source_type, source_id, position, style, metadata, created_at;

-- name: CreateGraphEdge :one
INSERT INTO graph_edges (
    graph_id,
    from_node_id,
    to_node_id,
    edge_type,
    is_directed,
    weight,
    confidence
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, graph_id, from_node_id, to_node_id, edge_type, is_directed, weight, confidence, style, metadata, created_at, updated_at;

-- name: DeleteGraphNodesByGraphID :exec
DELETE FROM graph_nodes WHERE graph_id = $1;