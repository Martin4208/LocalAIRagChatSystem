package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/google/uuid"
)

var (
	ErrDocumentNotFound  = errors.New("document not found")
	ErrAlreadyProcessing = errors.New("document is already being processed")
)

// ProcessOptions はドキュメント処理のオプション
type ProcessOptions struct {
	ChunkSize      int
	ChunkOverlap   int
	ForceReprocess bool
}

// DocumentProcessor はドキュメント処理のビジネスロジックを担当
type DocumentProcessor struct {
	queries      *db.Queries
	minioClient  *MinIOClient
	aiClient     client.AIWorkerClient
	qdrantClient client.QdrantClient
}

// NewDocumentProcessor は新しい DocumentProcessor を作成
func NewDocumentProcessor(
	queries *db.Queries,
	aiClient client.AIWorkerClient,
	qdrantClient client.QdrantClient,
) *DocumentProcessor {
	minioClient, err := NewMinIOClient(
		"localhost:9000",
		"admin",
		"password123",
		false,
	)
	if err != nil {
		log.Printf("Failed to create MinIO client: %v", err)
		return &DocumentProcessor{
			queries:     queries,
			minioClient: nil,
		}
	}

	return &DocumentProcessor{
		queries:      queries,
		minioClient:  minioClient,
		aiClient:     aiClient,
		qdrantClient: qdrantClient,
	}
}

// ProcessDocument はドキュメントを処理します
func (p *DocumentProcessor) ProcessDocument(
	ctx context.Context,
	workspaceID uuid.UUID,
	documentID uuid.UUID,
	opts ProcessOptions,
) error {
	doc, err := p.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:          documentID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return ErrDocumentNotFound
	}

	if doc.Status == "processing" && !opts.ForceReprocess {
		return ErrAlreadyProcessing
	}

	err = p.queries.UpdateDocumentStatus(ctx, db.UpdateDocumentStatusParams{
		Status:      "processing",
		ID:          documentID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	err = p.processDocumentInternal(ctx, doc, opts)
	if err != nil {
		updateErr := p.queries.UpdateDocumentStatus(ctx, db.UpdateDocumentStatusParams{
			Status:      "failed",
			ID:          documentID,
			WorkspaceID: workspaceID,
		})
		if updateErr != nil {
			log.Printf("Failed to update status to failed: %v", updateErr)
		}
		return err
	}

	return p.queries.UpdateDocumentStatus(ctx, db.UpdateDocumentStatusParams{
		Status:      "processed",
		ID:          documentID,
		WorkspaceID: workspaceID,
	})
}

// chunkWithPage はチャンクテキストとそのページ番号をセットで保持する内部型
type chunkWithPage struct {
	text       string
	pageNumber int
}

// processDocumentInternal は実際の処理を行います（内部用）
func (p *DocumentProcessor) processDocumentInternal(
	ctx context.Context,
	doc db.GetDocumentRow,
	opts ProcessOptions,
) error {
	log.Printf("Processing document: %s (type: %s)", doc.Name, doc.MimeType)

	if p.minioClient == nil {
		return fmt.Errorf("MinIO client not initialized")
	}

	// Step 1: MinIOからファイルをダウンロード
	reader, err := p.minioClient.DownloadFile(ctx, doc.MinioBucket, doc.MinioKey)
	if err != nil {
		return fmt.Errorf("failed to download file from MinIO: %w", err)
	}
	defer reader.Close()

	// Step 2: ファイルタイプに応じてチャンクを生成
	// ポイント：PDFはページ単位で処理し、page_numberを保持する
	var chunksWithPage []chunkWithPage
	chunker := NewTextChunker(opts.ChunkSize, opts.ChunkOverlap)

	switch doc.MimeType {
	case "application/pdf":
		// PDFはページ単位で抽出 → ページごとにチャンク化
		// これによりチャンクは必ず単一ページに属するようになる
		pages, err := NewPDFExtractor().ExtractPages(reader)
		if err != nil {
			return fmt.Errorf("failed to extract PDF pages: %w", err)
		}

		log.Printf("Extracted %d pages from PDF", len(pages))

		for _, page := range pages {
			// 1ページのテキストをチャンクに分割
			pageChunks := chunker.Chunk(page.Content)
			for _, chunk := range pageChunks {
				chunksWithPage = append(chunksWithPage, chunkWithPage{
					text:       chunk,
					pageNumber: page.PageNumber,
				})
			}
		}

	case "text/plain":
		// テキストファイルはページ概念がないのでpage_number=1固定
		data, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read text file: %w", err)
		}
		textChunks := chunker.Chunk(string(data))
		for _, chunk := range textChunks {
			chunksWithPage = append(chunksWithPage, chunkWithPage{
				text:       chunk,
				pageNumber: 1,
			})
		}

	default:
		return fmt.Errorf("unsupported file type: %s", doc.MimeType)
	}

	log.Printf("Split into %d chunks total", len(chunksWithPage))

	// Step 3: 既存チャンクを削除（再処理の場合）
	if opts.ForceReprocess {
		if err := p.queries.DeleteDocumentChunks(ctx, doc.ID); err != nil {
			log.Printf("Failed to delete existing chunks: %v", err)
		}
	}

	// Step 4: Embeddingを生成（テキストのみを渡す）
	texts := make([]string, len(chunksWithPage))
	for i, c := range chunksWithPage {
		texts[i] = c.text
	}

	log.Printf("Generating embeddings for %d chunks...", len(texts))
	embeddingResp, err := p.aiClient.EmbedDocuments(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	if len(embeddingResp.Embeddings) != len(chunksWithPage) {
		return fmt.Errorf("embedding count mismatch: got %d, expected %d",
			len(embeddingResp.Embeddings), len(chunksWithPage))
	}

	vectorDim := 1024
	if embeddingResp.Dim != nil {
		vectorDim = *embeddingResp.Dim
	}

	log.Printf("✅ Generated %d embeddings (dimension: %d)", embeddingResp.Count, vectorDim)

	// Step 5: チャンクをDBに保存（page_number付き）
	chunkIDs := make([]uuid.UUID, len(chunksWithPage))
	for i, c := range chunksWithPage {
		dbChunk, err := p.queries.CreateDocumentChunk(ctx, db.CreateDocumentChunkParams{
			DocumentID: doc.ID,
			ChunkIndex: int32(i),
			Content:    c.text,
			PageNumber: int32(c.pageNumber), // ← 追加：sqlcの再生成が必要
		})
		if err != nil {
			return fmt.Errorf("failed to save chunk %d: %w", i, err)
		}
		chunkIDs[i] = dbChunk.ID
	}

	// Step 6: QdrantにEmbeddingを保存（payloadにpage_numberを追加）
	log.Printf("Saving embeddings to Qdrant...")
	collectionName := fmt.Sprintf("workspace_%s", doc.WorkspaceID.String())

	if err := p.qdrantClient.EnsureCollection(ctx, collectionName, vectorDim); err != nil {
		return fmt.Errorf("failed to ensure Qdrant collection: %w", err)
	}

	points := make([]client.Point, len(chunksWithPage))
	for i, c := range chunksWithPage {
		points[i] = client.Point{
			ID:     chunkIDs[i].String(),
			Vector: embeddingResp.Embeddings[i],
			Payload: map[string]interface{}{
				"document_id":  doc.ID.String(),
				"workspace_id": doc.WorkspaceID.String(),
				"chunk_index":  i,
				"page_number":  c.pageNumber, // ← 追加：検索結果からページ番号を取得できるようになる
				"text":         c.text,
			},
		}
	}

	if err := p.qdrantClient.UpsertPoints(ctx, collectionName, points); err != nil {
		return fmt.Errorf("failed to save embeddings to Qdrant: %w", err)
	}

	log.Printf("✅ Saved %d embeddings to Qdrant collection '%s'", len(points), collectionName)
	log.Printf("Successfully processed document: %s (%d chunks)", doc.Name, len(chunksWithPage))

	return nil
}

// GetDocumentStatus はドキュメントの処理状況を取得します
func (p *DocumentProcessor) GetDocumentStatus(
	ctx context.Context,
	workspaceID uuid.UUID,
	documentID uuid.UUID,
) (map[string]interface{}, error) {
	doc, err := p.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:          documentID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, ErrDocumentNotFound
	}

	var processedAt interface{} = nil
	if doc.ProcessedAt.Valid {
		processedAt = doc.ProcessedAt.Time
	}

	chunkCount, err := p.queries.CountDocumentChunks(ctx, documentID)
	if err != nil {
		chunkCount = 0
	}

	status := map[string]interface{}{
		"status":       doc.Status,
		"progress":     nil,
		"error":        nil,
		"processed_at": processedAt,
	}

	if doc.Status == "processing" {
		status["progress"] = map[string]interface{}{
			"current_step":   "processing",
			"percentage":     50,
			"chunks_created": chunkCount,
		}
	}

	if doc.Status == "processed" {
		status["progress"] = map[string]interface{}{
			"chunks_created": chunkCount,
		}
	}

	return status, nil
}

// ChunkInfo はチャンク情報
type ChunkInfo struct {
	ID         uuid.UUID `json:"id"`
	ChunkIndex int       `json:"chunk_index"`
	PageNumber int       `json:"page_number"` // 追加
	Content    string    `json:"content"`
	CreatedAt  string    `json:"created_at"`
}

// GetDocumentChunks はドキュメントのチャンク一覧を取得します
func (p *DocumentProcessor) GetDocumentChunks(
	ctx context.Context,
	workspaceID uuid.UUID,
	documentID uuid.UUID,
	page int,
	limit int,
) ([]ChunkInfo, int, error) {
	_, err := p.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:          documentID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, 0, ErrDocumentNotFound
	}

	offset := (page - 1) * limit
	dbChunks, err := p.queries.GetDocumentChunks(ctx, db.GetDocumentChunksParams{
		DocumentID: documentID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get chunks: %w", err)
	}

	total, err := p.queries.CountDocumentChunks(ctx, documentID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chunks: %w", err)
	}

	chunks := make([]ChunkInfo, len(dbChunks))
	for i, dbChunk := range dbChunks {
		chunks[i] = ChunkInfo{
			ID:         dbChunk.ID,
			ChunkIndex: int(dbChunk.ChunkIndex),
			PageNumber: int(dbChunk.PageNumber), // sqlc再生成後に利用可能
			Content:    dbChunk.Content,
			CreatedAt:  dbChunk.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return chunks, int(total), nil
}
