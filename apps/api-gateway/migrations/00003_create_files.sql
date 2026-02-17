-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    sha256_hash TEXT NOT NULL UNIQUE,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    original_filename TEXT,
    minio_bucket TEXT NOT NULL,
    minio_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_accessed_at TIMESTAMPTZ
);
CREATE INDEX idx_files_sha256 ON files(sha256_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files CASCADE;
-- +goose StatementEnd
