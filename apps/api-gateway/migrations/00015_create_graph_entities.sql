-- +goose Up
-- +goose StatementBegin
CREATE TABLE graph_entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('person', 'organization', 'concept')),
    confidence FLOAT NOT NULL DEFAULT 0.8 CHECK(confidence >= 0 and confidence <= 1),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 基本的なB-treeインデックス
CREATE INDEX idx_entities_workspace ON graph_entities(workspace_id);

-- pg_trgm拡張（部分一致検索用）
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_entities_label_trgm ON graph_entities USING GIN(label gin_trgm_ops);

-- GINインデックス（JSONB用）
CREATE INDEX idx_entities_metadata ON graph_entities USING GIN(metadata);

-- 部分インデックス（論理削除フィルタ用）
CREATE INDEX idx_entities_active ON graph_entities(workspace_id) WHERE is_deleted = FALSE;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガーの設定
CREATE TRIGGER update_graph_entities_updated_at
    BEFORE UPDATE ON graph_entities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS graph_entities CASCADE;
-- +goose StatementEnd