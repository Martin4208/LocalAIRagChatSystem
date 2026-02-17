-- +goose Up
-- +goose StatementBegin
CREATE TABLE canvas_elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    canvas_id UUID NOT NULL REFERENCES canvases(id) ON DELETE CASCADE,
    element_type TEXT NOT NULL,
    position JSONB NOT NULL,
    z_index INTEGER NOT NULL DEFAULT 0,
    content JSONB NOT NULL,
    style JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_canvas_elements_canvas ON canvas_elements(canvas_id);
CREATE INDEX idx_canvas_elements_type ON canvas_elements(element_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS canvas_elements CASCADE;
-- +goose StatementEnd
