package service

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

// PageContent はPDFの1ページ分のテキストとページ番号を保持する。
// ここをstringの代わりに使うことで、チャンク生成時にページ情報が失われない。
type PageContent struct {
	PageNumber int    // 1-based（PDFの自然な表現に合わせる）
	Content    string // そのページから抽出したテキスト
}

// PDFExtractor はPDFからテキストを抽出します
type PDFExtractor struct{}

// NewPDFExtractor は新しい PDFExtractor を作成
func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

// ExtractPages はPDFをページ単位で抽出し、[]PageContent を返す。
// 旧: ExtractText (string) → 新: ExtractPages ([]PageContent)
// 変更理由：ページ境界情報をchunker/processorに渡すため。
func (e *PDFExtractor) ExtractPages(reader io.Reader) ([]PageContent, error) {
	// Step 1: データを全て読み込む
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF data: %w", err)
	}

	// Step 2: 一時ファイルに保存（pdfcpuがファイルパスを要求するため）
	tmpFile, err := os.CreateTemp("", "pdf-extract-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write PDF to temp file: %w", err)
	}
	tmpFile.Close()

	// Step 3: PDFコンテキストを作成
	ctx, err := api.ReadContextFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	// Step 4: ページ単位でテキストを抽出し、PageContentスライスに格納
	pageCount := ctx.PageCount
	pages := make([]PageContent, 0, pageCount)

	for i := 1; i <= pageCount; i++ {
		pageReader, err := pdfcpu.ExtractPageContent(ctx, i)
		if err != nil {
			// エラーページはスキップするが、ページ番号の連続性は維持する
			// （空ページとして記録しない：チャンクが生成されないだけでOK）
			continue
		}

		pageData, err := io.ReadAll(pageReader)
		if err != nil {
			continue
		}

		// クリーンアップしてから格納
		text := cleanExtractedText(string(pageData))
		if len(strings.TrimSpace(text)) == 0 {
			// 画像ページなどテキストが取れない場合はスキップ
			continue
		}

		pages = append(pages, PageContent{
			PageNumber: i,
			Content:    text,
		})
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no text found in PDF (possibly image-based or scanned document)")
	}

	return pages, nil
}

// ExtractText は後方互換性のために残す（テキストファイル処理などで引き続き使用）
// PDFには使わないこと。
func (e *PDFExtractor) ExtractText(reader io.Reader) (string, error) {
	pages, err := e.ExtractPages(reader)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, p := range pages {
		sb.WriteString(p.Content)
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String()), nil
}

// cleanExtractedText は抽出されたテキストをクリーンアップします
func cleanExtractedText(text string) string {
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(text)
}
