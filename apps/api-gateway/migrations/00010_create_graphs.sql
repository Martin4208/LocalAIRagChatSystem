-- +goose Up
-- +goose StatementBegin
CREATE TABLE graphs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    graph_type TEXT,
    layout_config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS graphs CASCADE;
-- +goose StatementEnd
