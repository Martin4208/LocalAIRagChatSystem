package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// AnalysisService は分析機能のビジネスロジックを担当
type AnalysisService struct {
	queries      *db.Queries
	aiClient     client.AIWorkerClient
	qdrantClient client.QdrantClient
	ollamaClient client.OllamaClient
}

// NewAnalysisService は新しいAnalysisServiceを作成
func NewAnalysisService(
	queries *db.Queries,
	aiClient client.AIWorkerClient,
	qdrantClient client.QdrantClient,
	ollamaClient client.OllamaClient,
) *AnalysisService {
	return &AnalysisService{
		queries:      queries,
		aiClient:     aiClient,
		qdrantClient: qdrantClient,
		ollamaClient: ollamaClient,
	}
}

// CreateAnalysis は分析ジョブを作成し、バックグラウンドで処理を開始
func (s *AnalysisService) CreateAnalysis(
	ctx context.Context,
	workspaceID uuid.UUID,
	title string,
	description *string,
	analysisType string,
	config map[string]interface{},
) (*db.Analysis, error) {
	// Step 1: config を JSONB に変換
	var configJSON pqtype.NullRawMessage
	if config != nil {
		configBytes, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		configJSON = pqtype.NullRawMessage{
			RawMessage: configBytes,
			Valid:      true,
		}
	}

	// Step 2: description を sql.NullString に変換
	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}

	// Step 3: DB に分析ジョブを作成（status = pending）
	analysis, err := s.queries.CreateAnalysis(ctx, db.CreateAnalysisParams{
		WorkspaceID:  workspaceID,
		Title:        title,
		Description:  desc,
		AnalysisType: analysisType,
		Config:       configJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create analysis: %w", err)
	}

	log.Printf("✅ Analysis created: %s (type: %s)", analysis.ID, analysisType)

	// Step 4: バックグラウンドで分析処理を開始（Goroutine）
	go func() {
		// 新しいコンテキストを作成（元のリクエストとは切り離す）
		bgCtx := context.Background()

		log.Printf("🔄 Starting background analysis: %s", analysis.ID)

		if err := s.ProcessAnalysis(bgCtx, workspaceID, analysis.ID); err != nil {
			log.Printf("❌ Analysis failed: %s - %v", analysis.ID, err)
		} else {
			log.Printf("✅ Analysis completed: %s", analysis.ID)
		}
	}()

	// Step 5: すぐにレスポンスを返す（202 Accepted）
	return &analysis, nil
}

// ProcessAnalysis は実際の分析処理を行います（バックグラウンド実行）
func (s *AnalysisService) ProcessAnalysis(
	ctx context.Context,
	workspaceID uuid.UUID,
	analysisID uuid.UUID,
) error {
	// Step 1: ステータスを processing に更新
	err := s.queries.UpdateAnalysisStatus(ctx, db.UpdateAnalysisStatusParams{
		ID:           analysisID,
		Status:       "processing",
		ErrorMessage: sql.NullString{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Step 2: 分析情報を取得
	analysis, err := s.queries.GetAnalysis(ctx, db.GetAnalysisParams{
		ID:          analysisID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		s.markAsFailed(ctx, analysisID, err)
		return fmt.Errorf("failed to get analysis: %w", err)
	}

	// Step 3: 対象ドキュメントを取得
	documents, err := s.getTargetDocuments(ctx, workspaceID, analysis.Config)
	if err != nil {
		s.markAsFailed(ctx, analysisID, err)
		return fmt.Errorf("failed to get documents: %w", err)
	}

	if len(documents) == 0 {
		err := fmt.Errorf("no documents found for analysis")
		s.markAsFailed(ctx, analysisID, err)
		return err
	}

	log.Printf("📄 Processing %d documents for analysis %s", len(documents), analysisID)

	// Step 4: 分析タイプに応じて処理分岐
	var results []db.CreateAnalysisResultParams
	switch analysis.AnalysisType {
	case "summary":
		results, err = s.processSummary(ctx, workspaceID, documents)
	case "keyword_extraction":
		results, err = s.processKeywordExtraction(ctx, workspaceID, documents)
	case "entity_recognition":
		results, err = s.processEntityRecognition(ctx, workspaceID, documents)
	default:
		err = fmt.Errorf("unsupported analysis type: %s", analysis.AnalysisType)
	}

	if err != nil {
		s.markAsFailed(ctx, analysisID, err)
		return fmt.Errorf("analysis processing failed: %w", err)
	}

	// Step 5: 結果を DB に保存
	for _, result := range results {
		result.AnalysisID = analysisID
		_, err := s.queries.CreateAnalysisResult(ctx, result)
		if err != nil {
			log.Printf("⚠️ Failed to save result: %v", err)
		}
	}

	// Step 6: ステータスを completed に更新
	err = s.queries.UpdateAnalysisStatus(ctx, db.UpdateAnalysisStatusParams{
		ID:           analysisID,
		Status:       "completed",
		ErrorMessage: sql.NullString{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to update status to completed: %w", err)
	}

	return nil
}

// getTargetDocuments は分析対象のドキュメントを取得
func (s *AnalysisService) getTargetDocuments(
	ctx context.Context,
	workspaceID uuid.UUID,
	config pqtype.NullRawMessage,
) ([]db.ListDocumentsRow, error) {
	// config から document_ids を取得
	var documentIDs []uuid.UUID

	if config.Valid && len(config.RawMessage) > 0 {
		var configMap map[string]interface{}
		if err := json.Unmarshal(config.RawMessage, &configMap); err == nil {
			if ids, ok := configMap["document_ids"].([]interface{}); ok {
				for _, id := range ids {
					if idStr, ok := id.(string); ok {
						if docID, err := uuid.Parse(idStr); err == nil {
							documentIDs = append(documentIDs, docID)
						}
					}
				}
			}
		}
	}

	// document_ids が指定されている場合は、それらを取得
	if len(documentIDs) > 0 {
		// TODO: 特定のドキュメントIDで絞り込む実装
		// 現在の ListDocuments は ID フィルタに非対応なので、全件取得後にフィルタ
		allDocs, err := s.queries.ListDocuments(ctx, db.ListDocumentsParams{
			WorkspaceID: workspaceID,
			Limit:       1000,
			Offset:      0,
		})
		if err != nil {
			return nil, err
		}

		// フィルタリング
		var filtered []db.ListDocumentsRow
		for _, doc := range allDocs {
			for _, targetID := range documentIDs {
				if doc.ID == targetID {
					filtered = append(filtered, doc)
					break
				}
			}
		}
		return filtered, nil
	}

	// document_ids が空の場合は、workspace内の全ドキュメントを対象
	return s.queries.ListDocuments(ctx, db.ListDocumentsParams{
		WorkspaceID: workspaceID,
		Limit:       1000,
		Offset:      0,
	})
}

// processSummary は要約分析を実行
func (s *AnalysisService) processSummary(
	ctx context.Context,
	workspaceID uuid.UUID,
	documents []db.ListDocumentsRow,
) ([]db.CreateAnalysisResultParams, error) {
	log.Printf("📝 Starting summary analysis for %d documents", len(documents))

	// Step 1: 全ドキュメントのチャンクを結合
	var allChunks []string
	maxChunks := 100 // 最大100チャンクまで（コンテキスト制限対策）

	for _, doc := range documents {
		chunks, err := s.queries.GetDocumentChunks(ctx, db.GetDocumentChunksParams{
			DocumentID: doc.ID,
			Limit:      int32(maxChunks),
			Offset:     0,
		})
		if err != nil {
			log.Printf("⚠️ Failed to get chunks for doc %s: %v", doc.ID, err)
			continue
		}

		for _, chunk := range chunks {
			allChunks = append(allChunks, chunk.Content)
			if len(allChunks) >= maxChunks {
				break
			}
		}

		if len(allChunks) >= maxChunks {
			break
		}
	}

	if len(allChunks) == 0 {
		return nil, fmt.Errorf("no chunks found in documents")
	}

	log.Printf("📊 Collected %d chunks for summarization", len(allChunks))

	// Step 2: チャンクを結合してプロンプト作成
	combinedText := strings.Join(allChunks, "\n\n---\n\n")

	prompt := fmt.Sprintf(`以下の資料群を分析し、日本語で要約を作成してください。

【要約の要件】
1. 主要なテーマを3-5個抽出してください
2. 各テーマについて2-3文で簡潔に説明してください
3. 重要なキーワードを太字で強調してください
4. 全体の結論を最後に1段落で述べてください

【資料内容】
%s

【要約】`, combinedText)

	// Step 3: LLM（Ollama）で要約生成
	log.Printf("🤖 Calling Ollama for summarization...")
	summary, err := s.ollamaClient.Generate(ctx, "phi3:mini", prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	log.Printf("✅ Summary generated: %d characters", len(summary))

	// Step 4: 結果を返す
	contentJSON, _ := json.Marshal(map[string]string{
		"summary": summary,
	})

	result := db.CreateAnalysisResultParams{
		ResultType:  "summary",
		Content:     contentJSON,
		ImageUrl:    sql.NullString{Valid: false},
		MinioBucket: sql.NullString{Valid: false},
		MinioKey:    sql.NullString{Valid: false},
		Metadata:    pqtype.NullRawMessage{Valid: false},
	}

	return []db.CreateAnalysisResultParams{result}, nil
}

// processKeywordExtraction はキーワード抽出分析を実行
func (s *AnalysisService) processKeywordExtraction(
	ctx context.Context,
	workspaceID uuid.UUID,
	documents []db.ListDocumentsRow,
) ([]db.CreateAnalysisResultParams, error) {
	// TODO: キーワード抽出の実装
	return nil, fmt.Errorf("keyword_extraction is not implemented yet")
}

// processEntityRecognition は固有表現抽出分析を実行
func (s *AnalysisService) processEntityRecognition(
	ctx context.Context,
	workspaceID uuid.UUID,
	documents []db.ListDocumentsRow,
) ([]db.CreateAnalysisResultParams, error) {
	// TODO: 固有表現抽出の実装
	return nil, fmt.Errorf("entity_recognition is not implemented yet")
}

// markAsFailed は分析を失敗としてマーク
func (s *AnalysisService) markAsFailed(ctx context.Context, analysisID uuid.UUID, err error) {
	updateErr := s.queries.UpdateAnalysisStatus(ctx, db.UpdateAnalysisStatusParams{
		ID:     analysisID,
		Status: "failed",
		ErrorMessage: sql.NullString{
			String: err.Error(),
			Valid:  true,
		},
	})
	if updateErr != nil {
		log.Printf("⚠️ Failed to mark analysis as failed: %v", updateErr)
	}
}

// jsonEscape は文字列をJSON用にエスケープ
func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// ListAnalyses は分析一覧を取得
func (s *AnalysisService) ListAnalyses(
	ctx context.Context,
	workspaceID uuid.UUID,
	status *string,
	limit int,
	offset int,
) ([]db.ListAnalysesRow, int64, error) {
	// ステータスフィルタの準備
	var statusFilter sql.NullString
	if status != nil {
		statusFilter = sql.NullString{String: *status, Valid: true}
	}

	// 一覧取得
	analyses, err := s.queries.ListAnalyses(ctx, db.ListAnalysesParams{
		WorkspaceID: workspaceID,
		Status:      statusFilter,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analyses: %w", err)
	}

	// 総数取得
	total, err := s.queries.CountAnalyses(ctx, db.CountAnalysesParams{
		WorkspaceID: workspaceID,
		Status:      statusFilter,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analyses: %w", err)
	}

	return analyses, total, nil
}

// GetAnalysis は分析詳細を取得
func (s *AnalysisService) GetAnalysis(
	ctx context.Context,
	workspaceID uuid.UUID,
	analysisID uuid.UUID,
) (*db.Analysis, error) {
	analysis, err := s.queries.GetAnalysis(ctx, db.GetAnalysisParams{
		ID:          analysisID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	return &analysis, nil
}

// GetAnalysisResults は分析結果を取得
func (s *AnalysisService) GetAnalysisResults(
	ctx context.Context,
	analysisID uuid.UUID,
) ([]db.AnalysisResult, error) {
	results, err := s.queries.ListAnalysisResults(ctx, analysisID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis results: %w", err)
	}

	return results, nil
}

// DeleteAnalysis は分析を削除（Soft delete）
func (s *AnalysisService) DeleteAnalysis(
	ctx context.Context,
	workspaceID uuid.UUID,
	analysisID uuid.UUID,
) error {
	err := s.queries.DeleteAnalysis(ctx, db.DeleteAnalysisParams{
		ID:          analysisID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete analysis: %w", err)
	}

	return nil
}
