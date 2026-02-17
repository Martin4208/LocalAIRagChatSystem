-- apps/api-gateway/queries/search.sql

-- name: GetChunksByIDs :many
-- page_number を追加して返す
SELECT 
    id,
    document_id,
    chunk_index,
    page_number,
    content,
    created_at
FROM document_chunks
WHERE id = ANY(@chunk_ids::uuid[]);

-- name: GetDocumentsByChunkIDs :many
SELECT DISTINCT
    d.id,
    d.name,
    d.workspace_id
FROM documents d
INNER JOIN document_chunks dc ON d.id = dc.document_id
WHERE dc.id = ANY(@chunk_ids::uuid[]);