-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_chunks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    qdrant_point_id UUID,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(document_id, chunk_index)
);
CREATE INDEX idx_chunks_document ON document_chunks(document_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS document_chunks CASCADE;
-- +goose StatementEnd
