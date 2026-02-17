-- +goose Up
-- +goose StatementBegin
CREATE TABLE analysis_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    result_type TEXT NOT NULL,
    content JSONB NOT NULL,
    image_url TEXT,
    minio_bucket TEXT,
    minio_key TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_analysis_results_analysis ON analysis_results(analysis_id);
CREATE INDEX idx_analysis_results_type ON analysis_results(result_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS analysis_results CASCADE;
-- +goose StatementEnd
