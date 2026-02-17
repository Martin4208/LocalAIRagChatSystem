package service

import (
	"strings"
)

// TextChunker はテキストをチャンクに分割します
type TextChunker struct {
	ChunkSize    int
	ChunkOverlap int
}

// NewTextChunker は新しい TextChunker を作成
func NewTextChunker(chunkSize, chunkOverlap int) *TextChunker {
	// バリデーション
	if chunkSize <= 0 {
		chunkSize = 500
	}
	if chunkOverlap < 0 {
		chunkOverlap = 0
	}
	if chunkOverlap >= chunkSize {
		chunkOverlap = chunkSize / 10 // チャンクサイズの10%
	}

	return &TextChunker{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

// Chunk はテキストをチャンクに分割します（文の境界を考慮）
func (c *TextChunker) Chunk(text string) []string {
	chunks := []string{}

	// テキストが空の場合
	if len(strings.TrimSpace(text)) == 0 {
		return chunks
	}

	// rune（文字）に変換（日本語対応）
	runes := []rune(text)
	textLen := len(runes)

	// ChunkSizeより短い場合はそのまま返す
	if textLen <= c.ChunkSize {
		return []string{text}
	}

	start := 0
	for start < textLen {
		// チャンクの終了位置を計算
		end := start + c.ChunkSize
		if end > textLen {
			end = textLen
		}

		// 文の境界を探す（最後のチャンク以外）
		if end < textLen {
			end = c.findSentenceBoundary(runes, start, end)
		}

		// チャンクを切り出す
		chunk := string(runes[start:end])
		chunk = strings.TrimSpace(chunk)

		// 空でないチャンクのみ追加
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}

		// 次のチャンクの開始位置（オーバーラップ分を考慮）
		nextStart := end - c.ChunkOverlap

		// オーバーラップが大きすぎて進まない場合の対策
		if nextStart <= start {
			nextStart = end
		}

		start = nextStart
	}

	return chunks
}

// findSentenceBoundary は文の境界を探します
func (c *TextChunker) findSentenceBoundary(runes []rune, start, idealEnd int) int {
	// idealEnd の前後20文字以内で句点を探す
	searchStart := idealEnd - 20
	if searchStart < start {
		searchStart = start
	}

	searchEnd := idealEnd + 20
	if searchEnd > len(runes) {
		searchEnd = len(runes)
	}

	// まず idealEnd から後ろを探す（できるだけ長くする）
	for i := idealEnd; i < searchEnd; i++ {
		if c.isSentenceEnd(runes, i) {
			return i + 1 // 句点の次の文字から
		}
	}

	// 見つからなければ idealEnd から前を探す
	for i := idealEnd - 1; i >= searchStart; i-- {
		if c.isSentenceEnd(runes, i) {
			return i + 1
		}
	}

	// どちらも見つからなければ idealEnd を返す
	return idealEnd
}

// isSentenceEnd は文末判定を行います
func (c *TextChunker) isSentenceEnd(runes []rune, pos int) bool {
	if pos >= len(runes) {
		return false
	}

	r := runes[pos]

	// 日本語の句点
	if r == '。' || r == '！' || r == '？' {
		return true
	}

	// 英語の句点（後ろにスペースがある場合）
	if (r == '.' || r == '!' || r == '?') && pos+1 < len(runes) {
		next := runes[pos+1]
		if next == ' ' || next == '\n' {
			return true
		}
	}

	// 改行（段落の境界）
	if r == '\n' && pos+1 < len(runes) {
		next := runes[pos+1]
		if next == '\n' {
			return true // 連続する改行は段落の境界
		}
	}

	return false
}

// ChunkSimple はシンプルな固定文字数分割（文の境界を考慮しない）
func (c *TextChunker) ChunkSimple(text string) []string {
	chunks := []string{}

	if len(strings.TrimSpace(text)) == 0 {
		return chunks
	}

	runes := []rune(text)
	textLen := len(runes)

	if textLen <= c.ChunkSize {
		return []string{text}
	}

	start := 0
	for start < textLen {
		end := start + c.ChunkSize
		if end > textLen {
			end = textLen
		}

		chunk := string(runes[start:end])
		chunk = strings.TrimSpace(chunk)

		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}

		start = end - c.ChunkOverlap
		if start <= 0 {
			start = end
		}
	}

	return chunks
}

// GetChunkCount はテキストを分割した場合のチャンク数を推定します
func (c *TextChunker) GetChunkCount(text string) int {
	textLen := len([]rune(text))
	if textLen <= c.ChunkSize {
		return 1
	}

	step := c.ChunkSize - c.ChunkOverlap
	if step <= 0 {
		step = c.ChunkSize
	}

	return (textLen + step - 1) / step
}
