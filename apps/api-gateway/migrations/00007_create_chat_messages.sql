-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    message_index INTEGER NOT NULL,
    document_refs JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(chat_id, message_index)
);
CREATE INDEX idx_messages_chat ON chat_messages(chat_id, message_index);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat_messages CASCADE;
-- +goose StatementEnd
