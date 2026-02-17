package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/google/uuid"
)

// ChatService はチャット機能のビジネスロジックを担当
type ChatService struct {
	queries        *db.Queries
	aiWorkerClient client.AIWorkerClient
	qdrantClient   client.QdrantClient
	ollamaClient   client.OllamaClient
}

// NewChatService は新しいChatServiceを作成
func NewChatService(
	queries *db.Queries,
	aiWorkerClient client.AIWorkerClient,
	qdrantClient client.QdrantClient,
	ollamaClient client.OllamaClient,
) *ChatService {
	return &ChatService{
		queries:        queries,
		aiWorkerClient: aiWorkerClient,
		qdrantClient:   qdrantClient,
		ollamaClient:   ollamaClient,
	}
}

// GenerateResponse はRAGを使ってLLMの応答を生成します
func (s *ChatService) GenerateResponse(
	ctx context.Context,
	workspaceID uuid.UUID,
	chatID uuid.UUID,
	userMessage string,
) (string, []api.DocumentReference, error) {
	log.Printf("🔍 [RAG] Starting response generation for message: %s", userMessage)

	// Step 1: ユーザーメッセージをEmbedding化
	embedResp, err := s.aiWorkerClient.EmbedDocuments(ctx, []string{userMessage})
	if err != nil {
		return "", nil, fmt.Errorf("failed to embed user message: %w", err)
	}

	if len(embedResp.Embeddings) == 0 {
		return "", nil, fmt.Errorf("no embeddings returned from AI worker")
	}

	queryVector := embedResp.Embeddings[0]

	// Step 2: Qdrantで類似チャンクを検索
	collectionName := fmt.Sprintf("workspace_%s", workspaceID.String())
	searchResp, err := s.qdrantClient.Search(ctx, collectionName, queryVector, 5)
	if err != nil {
		return "", nil, fmt.Errorf("failed to search Qdrant: %w", err)
	}
	log.Printf("✅ [RAG] Found %d results from Qdrant", len(searchResp.Result))

	// Step 3: 検索結果からDocumentReferenceを生成（page_number付き）
	documentRefs := s.extractDocumentRefs(searchResp.Result)

	// Step 4: コンテキストを構築
	context := s.buildContext(searchResp.Result)

	// Step 5: プロンプトを作成してOllamaで生成
	prompt := s.buildPrompt(context, userMessage)
	llmResponse, err := s.ollamaClient.Generate(ctx, "phi3:mini", prompt)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate response from LLM: %w", err)
	}

	log.Printf("✅ [RAG] Response generated: %d characters", len(llmResponse))

	return llmResponse, documentRefs, nil
}

// extractDocumentRefs は Qdrant の検索結果から DocumentReference スライスを生成する。
// payload に page_number が含まれるようになったので、それを取り出してセットする。
func (s *ChatService) extractDocumentRefs(
	results []client.SearchResult,
) []api.DocumentReference {
	refs := make([]api.DocumentReference, 0, len(results))

	for _, r := range results {
		payload := r.Payload

		// document_id
		docIDStr, ok := payload["document_id"].(string)
		if !ok {
			continue
		}
		docUUID, err := uuid.Parse(docIDStr)
		if err != nil {
			continue
		}

		// chunk_index
		chunkIndex := 0
		if v, ok := payload["chunk_index"].(float64); ok {
			chunkIndex = int(v)
		}

		// page_number（ない場合はnilのまま：後方互換）
		var pageNumber *int
		if v, ok := payload["page_number"].(float64); ok {
			pn := int(v)
			pageNumber = &pn
		}

		// content_preview（最初の200文字）
		contentPreview := ""
		if v, ok := payload["text"].(string); ok {
			runes := []rune(v)
			if len(runes) > 200 {
				contentPreview = string(runes[:200])
			} else {
				contentPreview = v
			}
		}

		score := float32(r.Score)
		ref := api.DocumentReference{
			DocumentId:     docUUID,
			ChunkIndex:     int32(chunkIndex),
			PageNumber:     pageNumber, // ← 追加（openapi再生成後に型が確定する）
			Score:          score,
			ContentPreview: &contentPreview,
		}

		refs = append(refs, ref)
	}

	return refs
}

// buildContext は検索結果からコンテキストテキストを構築
func (s *ChatService) buildContext(results []client.SearchResult) string {
	var contextParts []string

	for i, result := range results {
		content, ok := result.Payload["text"].(string)
		if !ok {
			continue
		}

		// ページ番号が取れる場合はコンテキストにも含める（LLMへのヒントになる）
		pageInfo := ""
		if v, ok := result.Payload["page_number"].(float64); ok {
			pageInfo = fmt.Sprintf(" [P.%d]", int(v))
		}

		contextParts = append(contextParts, fmt.Sprintf(
			"--- Document %d%s (Score: %.3f) ---\n%s",
			i+1, pageInfo, result.Score, content,
		))
	}

	if len(contextParts) == 0 {
		return "関連する資料が見つかりませんでした。"
	}

	return strings.Join(contextParts, "\n\n")
}

// buildPrompt はLLMに送るプロンプトを構築
func (s *ChatService) buildPrompt(context string, userMessage string) string {
	systemPrompt := `あなたは提供された資料を基に正確に回答するAIアシスタントです。
以下のルールに従ってください：
1. 提供された資料の内容のみを基に回答する
2. 資料に記載がない場合は「資料には記載されていません」と答える
3. 推測や一般知識での回答は避ける
4. 回答は簡潔かつ正確に`

	return fmt.Sprintf(`%s

参考資料:
%s

ユーザーの質問: %s

回答:`, systemPrompt, context, userMessage)
}
