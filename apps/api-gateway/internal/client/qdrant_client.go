package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// QdrantClient はQdrantとの通信インターフェース
type QdrantClient interface {
	// EnsureCollection はCollectionが存在しなければ作成します
	EnsureCollection(ctx context.Context, collectionName string, vectorSize int) error

	// UpsertPoints はEmbeddingベクトルを保存します
	UpsertPoints(ctx context.Context, collectionName string, points []Point) error

	// Search は類似検索を実行します
	Search(ctx context.Context, collectionName string, vector []float64, limit int) (*SearchResponse, error)

	// document_idでポイントを削除
	DeletePointsByDocumentID(ctx context.Context, collectionName string, documentID string) error
}

// qdrantClient はQdrantClientの実装
type qdrantClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewQdrantClient は新しいQdrant Clientを作成します
func NewQdrantClient(baseURL string) QdrantClient {
	return &qdrantClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// EnsureCollection はCollectionが存在しなければ作成します
func (c *qdrantClient) EnsureCollection(ctx context.Context, collectionName string, vectorSize int) error {
	// Step 1: Collectionが存在するかチェック
	checkURL := fmt.Sprintf("%s/collections/%s", c.baseURL, collectionName)

	req, err := http.NewRequestWithContext(ctx, "GET", checkURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}
	defer resp.Body.Close()

	// Step 2: 既に存在する場合は何もしない
	if resp.StatusCode == http.StatusOK {
		return nil // 既に存在する
	}

	// Step 3: 存在しない場合は作成
	createReq := CreateCollectionRequest{
		Vectors: VectorConfig{
			Size:     vectorSize,
			Distance: "Cosine", // コサイン類似度
		},
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	createURL := fmt.Sprintf("%s/collections/%s", c.baseURL, collectionName)
	req, err = http.NewRequestWithContext(
		ctx,
		"PUT",
		createURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create collection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: ステータスコードチェック
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create collection: status %d", resp.StatusCode)
	}

	return nil
}

// UpsertPoints はEmbeddingベクトルを保存します
func (c *qdrantClient) UpsertPoints(ctx context.Context, collectionName string, points []Point) error {
	// Step 1: リクエストボディ作成
	reqBody := UpsertPointsRequest{
		Points: points,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal points: %w", err)
	}

	// Step 2: PUT リクエスト作成
	url := fmt.Sprintf("%s/collections/%s/points", c.baseURL, collectionName)

	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upsert points: %w", err)
	}
	defer resp.Body.Close()

	// Step 3: ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upsert points: status %d", resp.StatusCode)
	}

	return nil
}

// Search は類似検索を実行します
func (c *qdrantClient) Search(ctx context.Context, collectionName string, vector []float64, limit int) (*SearchResponse, error) {
	// Step 1: リクエストボディ作成
	reqBody := SearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true, // Payload（メタデータ）も取得
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Step 2: POST リクエスト作成
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

	// Step 3: リクエスト送信
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: status %d", resp.StatusCode)
	}

	// Step 5: レスポンスデコード
	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// DeletePointsByDocumentID はdocument_idでポイントを削除します
func (c *qdrantClient) DeletePointsByDocumentID(
	ctx context.Context,
	collectionName string,
	documentID string,
) error {
	// Step 1: リクエストボディ作成
	reqBody := map[string]interface{}{
		"filter": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"key": "document_id",
					"match": map[string]interface{}{
						"value": documentID,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal delete request: %w", err)
	}

	// Step 2: POST リクエスト作成
	url := fmt.Sprintf("%s/collections/%s/points/delete", c.baseURL, collectionName)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Step 3: リクエスト送信
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete points: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete points: status %d", resp.StatusCode)
	}

	return nil
}
