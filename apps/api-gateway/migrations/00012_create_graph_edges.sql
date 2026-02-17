-- +goose Up
-- +goose StatementBegin
CREATE TABLE graph_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    graph_id UUID NOT NULL REFERENCES graphs(id) ON DELETE CASCADE,
    from_node_id UUID NOT NULL REFERENCES graph_nodes(id) ON DELETE CASCADE,
    to_node_id UUID NOT NULL REFERENCES graph_nodes(id) ON DELETE CASCADE, 
    edge_type TEXT NOT NULL,
    is_directed BOOLEAN NOT NULL DEFAULT true,
    weight FLOAT,
    confidence FLOAT,
    style JSONB,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(graph_id, from_node_id, to_node_id, edge_type)
);
CREATE INDEX idx_graph_edges_graph ON graph_edges(graph_id);
CREATE INDEX idx_graph_edges_from ON graph_edges(from_node_id);
CREATE INDEX idx_graph_edges_to ON graph_edges(to_node_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS graph_edges CASCADE;
-- +goose StatementEnd
