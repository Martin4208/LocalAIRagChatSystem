-- name: GetDocumentMetadataByIDs :many
-- ドキュメントのメタデータのみ取得（チャンク内容は不要）
SELECT 
    d.id,
    d.name,
    f.mime_type,
    f.size_bytes,
    d.created_at,
    d.updated_at
FROM documents d
INNER JOIN files f ON d.file_id = f.id
WHERE 
    d.id = ANY(@document_ids::uuid[])
    AND d.workspace_id = @workspace_id
    AND d.deleted_at IS NULL
ORDER BY d.name;