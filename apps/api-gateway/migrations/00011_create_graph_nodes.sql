-- +goose Up
-- +goose StatementBegin
CREATE TABLE graph_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    graph_id UUID  NOT NULL REFERENCES graphs(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    node_type TEXT NOT NULL,
    source_type TEXT,
    source_id UUID,
    position JSONB,
    style JSONB,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_source FOREIGN KEY (source_id) REFERENCES document_chunks(id) ON DELETE SET NULL
);
CREATE INDEX idx_graph_nodes_graph ON graph_nodes(graph_id);
CREATE INDEX idx_graph_nodes_source ON graph_nodes(source_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS graph_nodes CASCADE;
-- +goose StatementEnd
