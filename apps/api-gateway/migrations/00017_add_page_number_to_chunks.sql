-- +goose Up
-- +goose StatementBegin

-- page_number を専用カラムとして追加する（JSONB ではなく専用カラムにする理由：
-- インデックスを張れる・型安全・将来的なページフィルタクエリが高速になる）
ALTER TABLE document_chunks ADD COLUMN page_number INTEGER NOT NULL DEFAULT 0;

-- ページ番号での検索・フィルタに備えてインデックスを追加
CREATE INDEX idx_chunks_page_number ON document_chunks(document_id, page_number);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chunks_page_number;
ALTER TABLE document_chunks DROP COLUMN IF EXISTS page_number;
-- +goose StatementEnd