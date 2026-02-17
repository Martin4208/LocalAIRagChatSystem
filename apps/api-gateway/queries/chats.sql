-- ========================================
-- Chat Operations
-- ========================================

-- name: CreateChat :one
INSERT INTO chats (
    workspace_id,
    title,
    filter_config,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, now(), now()
)
RETURNING *;

-- name: ListChats :many
SELECT 
    c.id,
    c.workspace_id,
    c.title,
    c.filter_config,
    c.created_at,
    c.updated_at
FROM chats c
WHERE 
    c.workspace_id = $1 
    AND c.deleted_at IS NULL
ORDER BY c.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CountChats :one
SELECT COUNT(*) 
FROM chats 
WHERE 
    workspace_id = $1 
    AND deleted_at IS NULL;

-- name: GetChat :one
SELECT * FROM chats
WHERE 
    id = $1 
    AND workspace_id = $2 
    AND deleted_at IS NULL;

-- name: DeleteChat :exec
UPDATE chats
SET 
    deleted_at = now(),
    updated_at = now()
WHERE 
    id = $1 
    AND workspace_id = $2 
    AND deleted_at IS NULL;

-- name: UpdateChatTimestamp :exec
UPDATE chats
SET updated_at = now()
WHERE id = $1;

-- name: GetChatMessageCount :one
SELECT COUNT(*)::bigint as count
FROM chat_messages
WHERE chat_id = $1;

-- name: GetLastMessageTime :one
SELECT 
    CASE 
        WHEN COUNT(*) > 0 THEN MAX(created_at)
        ELSE NULL
    END::timestamptz as last_time
FROM chat_messages
WHERE chat_id = $1;

-- ========================================
-- Chat Message Operations
-- ========================================

-- name: GetChatMessages :many
SELECT 
    id,
    chat_id,
    role,
    content,
    message_index,
    document_refs,
    created_at
FROM chat_messages
WHERE 
    chat_id = $1
ORDER BY message_index ASC
LIMIT $2;

-- name: GetMaxMessageIndex :one
SELECT COALESCE(MAX(message_index), -1)::int as max_index
FROM chat_messages
WHERE chat_id = $1;

-- name: CreateChatMessage :one
INSERT INTO chat_messages (
    chat_id,
    role,
    content,
    message_index,
    document_refs,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, now()
)
RETURNING *;