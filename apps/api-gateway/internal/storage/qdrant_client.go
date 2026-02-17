// apps/api-gateway/internal/storage/qdrant_client.go

package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// QdrantClient は Qdrant との通信を担当
type QdrantClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewQdrantClient は新しい QdrantClient を作成
func NewQdrantClient(baseURL string) *QdrantClient {
	return &QdrantClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchRequest は Qdrant 検索リクエスト
type SearchRequest struct {
	Vector      []float64 `json:"vector"`
	Limit       int       `json:"limit"`
	WithPayload bool      `json:"with_payload"`
}

// SearchResult は検索結果の1件
type SearchResult struct {
	ID      string                 `json:"id"`
	Score   float64                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

// SearchResponse は Qdrant からの検索レスポンス
type SearchResponse struct {
	Result []SearchResult `json:"result"`
}

// Search はベクトル類似検索を実行
func (c *QdrantClient) Search(
	ctx context.Context,
	collectionName string,
	vector []float64,
	limit int,
) (*SearchResponse, error) {
	// リクエストボディ作成
	reqBody := SearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// POST リクエスト作成
	url := fmt.Sprintf("%s/collections/%s/points/search", c.baseURL, collectionName)
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// リクエスト送信
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: status %d", resp.StatusCode)
	}

	// レスポンスデコード
	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
