-- +goose Up
-- +goose StatementBegin
CREATE TABLE graph_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    source_entity_id UUID NOT NULL REFERENCES graph_entities(id) ON DELETE CASCADE,
    target_entity_id UUID NOT NULL REFERENCES graph_entities(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    is_directed BOOLEAN NOT NULL DEFAULT TRUE,
    weight FLOAT NOT NULL DEFAULT 1.0,
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (source_entity_id != target_entity_id),
    CHECK (weight > 0),
    CHECK (
        (valid_from IS NULL AND valid_to IS NULL) OR
        (valid_from IS NOT NULL AND valid_to IS NOT NULL AND valid_from < valid_to) OR
        (valid_from IS NOT NULL AND valid_to IS NULL)
    )
);

-- インデックス
CREATE INDEX idx_relations_workspace ON graph_relations(workspace_id);
CREATE INDEX idx_relations_source ON graph_relations(source_entity_id);
CREATE INDEX idx_relations_target ON graph_relations(target_entity_id);
CREATE INDEX idx_relations_type ON graph_relations(type);
CREATE INDEX idx_relations_active ON graph_relations(workspace_id) WHERE is_deleted = FALSE;

-- updated_atトリガー（関数は既に存在）
CREATE TRIGGER update_graph_relations_updated_at
    BEFORE UPDATE ON graph_relations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS graph_relations CASCADE;
-- +goose StatementEnd