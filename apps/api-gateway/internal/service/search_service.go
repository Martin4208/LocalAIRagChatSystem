// apps/api-gateway/internal/service/search_service.go

package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/google/uuid"
)

// SearchService は RAG 検索のビジネスロジック
type SearchService struct {
	queries  *db.Queries
	aiWorker client.AIWorkerClient
	qdrant   client.QdrantClient
}

// NewSearchService は新しい SearchService を作成
func NewSearchService(
	queries *db.Queries,
	aiWorker client.AIWorkerClient,
	qdrant client.QdrantClient,
) *SearchService {
	return &SearchService{
		queries:  queries,
		aiWorker: aiWorker,
		qdrant:   qdrant,
	}
}

// SearchResult は検索結果
type SearchResult struct {
	Answer  string         `json:"answer"`
	Sources []SearchSource `json:"sources"`
}

// SearchSource は検索ソース情報
type SearchSource struct {
	DocumentID uuid.UUID `json:"document_id"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	Score      float64   `json:"score"`
}

// Search は RAG 検索を実行
func (s *SearchService) Search(
	ctx context.Context,
	workspaceID uuid.UUID,
	query string,
	topK int,
) (*SearchResult, error) {
	log.Printf("Starting RAG search: workspace=%s, query=%s, topK=%d", workspaceID, query, topK)

	// コレクション名を動的に生成
	collectionName := fmt.Sprintf("workspace_%s", workspaceID)
	log.Printf("Using Qdrant collection: %s", collectionName)

	// Step 1: クエリを Embedding 化
	log.Printf("[1/4] Embedding query...")
	queryVector, err := s.aiWorker.EmbedQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	log.Printf("Query embedded: dim=%d", len(queryVector))

	// Step 2: Qdrant で類似検索
	log.Printf("[2/4] Searching Qdrant...")
	searchResp, err := s.qdrant.Search(ctx, collectionName, queryVector, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search Qdrant: %w", err)
	}
	log.Printf("Found %d results from Qdrant", len(searchResp.Result))

	if len(searchResp.Result) == 0 {
		return &SearchResult{
			Answer:  "申し訳ありません。関連する情報が見つかりませんでした。",
			Sources: []SearchSource{},
		}, nil
	}

	// Step 3: PostgreSQL から元テキストを取得
	log.Printf("[3/4] Fetching chunks from PostgreSQL...")
	chunkIDs := make([]uuid.UUID, len(searchResp.Result))
	scoreMap := make(map[uuid.UUID]float64)

	for i, result := range searchResp.Result {
		chunkID, err := uuid.Parse(result.ID)
		if err != nil {
			log.Printf("Warning: invalid chunk ID: %s", result.ID)
			continue
		}
		chunkIDs[i] = chunkID
		scoreMap[chunkID] = result.Score
	}

	chunks, err := s.queries.GetChunksByIDs(ctx, chunkIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks: %w", err)
	}
	log.Printf("Retrieved %d chunks from PostgreSQL", len(chunks))

	// Step 4: コンテキストを作成して LLM に投げる
	log.Printf("[4/4] Generating answer with LLM...")
	contextTexts := make([]string, len(chunks))
	sources := make([]SearchSource, len(chunks))

	for i, chunk := range chunks {
		contextTexts[i] = chunk.Content
		sources[i] = SearchSource{
			DocumentID: chunk.DocumentID,
			ChunkIndex: int(chunk.ChunkIndex),
			Content:    chunk.Content,
			Score:      scoreMap[chunk.ID],
		}
	}

	// LLM で回答生成
	generateResp, err := s.aiWorker.GenerateAnswer(ctx, query, contextTexts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate answer: %w", err)
	}
	log.Printf("Answer generated successfully")

	return &SearchResult{
		Answer:  generateResp.Answer,
		Sources: sources,
	}, nil
}
