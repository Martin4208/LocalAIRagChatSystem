package service

import (
	"strings"
	"testing"
)

func TestTextChunker_Chunk_ShortText(t *testing.T) {
	chunker := NewTextChunker(500, 50)
	text := "これは短いテキストです。"

	chunks := chunker.Chunk(text)

	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk, got %d", len(chunks))
	}

	if chunks[0] != text {
		t.Errorf("Chunk content mismatch")
	}
}

func TestTextChunker_Chunk_ExactSize(t *testing.T) {
	chunker := NewTextChunker(500, 50)
	text := strings.Repeat("あ", 500)

	chunks := chunker.Chunk(text)

	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk, got %d", len(chunks))
	}
}

func TestTextChunker_Chunk_LongText(t *testing.T) {
	chunker := NewTextChunker(100, 20)

	// 300文字のテキスト
	text := strings.Repeat("これはテストです。", 30)

	chunks := chunker.Chunk(text)

	// 少なくとも2チャンクは期待
	if len(chunks) < 2 {
		t.Errorf("Expected at least 2 chunks, got %d", len(chunks))
	}

	t.Logf("Generated %d chunks from %d characters", len(chunks), len([]rune(text)))

	// 各チャンクのサイズをチェック
	for i, chunk := range chunks {
		runeCount := len([]rune(chunk))
		t.Logf("Chunk %d: %d characters", i+1, runeCount)

		// 最後のチャンク以外はサイズ制限を超えないこと
		if i < len(chunks)-1 && runeCount > chunker.ChunkSize+50 {
			t.Errorf("Chunk %d is too large: %d characters", i+1, runeCount)
		}
	}
}

func TestTextChunker_Chunk_WithSentences(t *testing.T) {
	chunker := NewTextChunker(50, 10)

	text := "これは最初の文です。これは二番目の文です。これは三番目の文です。これは四番目の文です。これは五番目の文です。"

	chunks := chunker.Chunk(text)

	t.Logf("Generated %d chunks", len(chunks))

	for i, chunk := range chunks {
		t.Logf("Chunk %d: %s", i+1, chunk)
	}

	// 少なくとも2チャンク
	if len(chunks) < 2 {
		t.Errorf("Expected at least 2 chunks, got %d", len(chunks))
	}
}

func TestTextChunker_ChunkSimple(t *testing.T) {
	chunker := NewTextChunker(100, 20)

	text := strings.Repeat("あいうえお", 50) // 250文字

	chunks := chunker.ChunkSimple(text)

	t.Logf("Simple chunking: %d chunks from %d characters", len(chunks), len([]rune(text)))

	if len(chunks) < 2 {
		t.Errorf("Expected at least 2 chunks, got %d", len(chunks))
	}
}

func TestTextChunker_GetChunkCount(t *testing.T) {
	chunker := NewTextChunker(500, 50)

	tests := []struct {
		textLen  int
		expected int
	}{
		{100, 1},
		{500, 1},
		{501, 2},
		{1000, 3},
		{1500, 4},
	}

	for _, tt := range tests {
		text := strings.Repeat("あ", tt.textLen)
		count := chunker.GetChunkCount(text)

		if count < tt.expected-1 || count > tt.expected+1 {
			t.Errorf("Text length %d: expected ~%d chunks, got %d", tt.textLen, tt.expected, count)
		}
	}
}

func TestTextChunker_EmptyText(t *testing.T) {
	chunker := NewTextChunker(500, 50)

	chunks := chunker.Chunk("")
	if len(chunks) != 0 {
		t.Errorf("Expected 0 chunks for empty text, got %d", len(chunks))
	}

	chunks = chunker.Chunk("   \n\n  ")
	if len(chunks) != 0 {
		t.Errorf("Expected 0 chunks for whitespace-only text, got %d", len(chunks))
	}
}

func TestTextChunker_Validation(t *testing.T) {
	// 無効なパラメータでも動作すること
	chunker := NewTextChunker(-1, -1)

	if chunker.ChunkSize <= 0 {
		t.Error("ChunkSize should be positive")
	}

	if chunker.ChunkOverlap < 0 {
		t.Error("ChunkOverlap should be non-negative")
	}
}
