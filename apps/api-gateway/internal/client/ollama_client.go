package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// OllamaClient はOllama APIとの通信インターフェース
type OllamaClient interface {
	// Generate はLLMを使ってテキスト生成を行います
	Generate(ctx context.Context, model string, prompt string) (string, error)
	WarmUp(ctx context.Context, model string) error
}

// ollamaClient はOllamaClientの実装
type ollamaClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewOllamaClient は新しいOllama Clientを作成します
func NewOllamaClient(baseURL string) OllamaClient {
	return &ollamaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Minute, // LLM生成は時間がかかる可能性があるため長めに設定
		},
	}
}

// OllamaGenerateRequest はOllama APIのリクエスト形式
type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaGenerateResponse はOllama APIのレスポンス形式（Non-streaming）
type OllamaGenerateResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// WarmUp はアプリ起動時にOllamaを準備状態にします
func (c *ollamaClient) WarmUp(ctx context.Context, model string) error {
	log.Println("🔥 [Ollama] Starting warmup...")

	// Step 1: Health check
	if err := c.waitForReady(ctx, 60*time.Second); err != nil {
		return fmt.Errorf("ollama not ready: %w", err)
	}

	// Step 2: モデルが存在するか確認
	exists, err := c.checkModelExists(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to check model: %w", err)
	}

	if !exists {
		log.Printf("⚠️  [Ollama] Model '%s' not found", model)

		// 環境変数で自動プルを制御
		if os.Getenv("OLLAMA_AUTO_PULL") != "true" {
			return fmt.Errorf("model '%s' not found. Run 'ollama pull %s' first, or set OLLAMA_AUTO_PULL=true", model, model)
		}

		log.Printf("📥 [Ollama] Auto-pulling model '%s'...", model)
		if err := c.pullModel(ctx, model); err != nil {
			return fmt.Errorf("failed to pull model: %w", err)
		}
	} else {
		log.Printf("✅ [Ollama] Model '%s' already exists", model)
	}

	// Step 3: ダミーリクエストでモデルをメモリにロード
	log.Printf("🔄 [Ollama] Loading model into memory: %s", model)

	_, err = c.Generate(ctx, model, "warmup")
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}

	log.Printf("✅ [Ollama] Model loaded: %s", model)
	log.Println("✅ [Ollama] Warmup complete!")
	return nil
}

// checkModelExists はモデルが存在するかチェック
func (c *ollamaClient) checkModelExists(ctx context.Context, model string) (bool, error) {
	url := fmt.Sprintf("%s/api/tags", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	for _, m := range result.Models {
		if m.Name == model {
			return true, nil
		}
	}

	return false, nil
}

// pullModel はモデルをダウンロード
// pullModel はモデルをダウンロード
func (c *ollamaClient) pullModel(ctx context.Context, model string) error {
	log.Printf("📥 [Ollama] Pulling model '%s' (this may take 2-5 minutes)...", model)

	url := fmt.Sprintf("%s/api/pull", c.baseURL)

	reqBody := map[string]interface{}{
		"name": model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pull failed with status %d: %s", resp.StatusCode, string(body))
	}

	// プル進捗をストリーミングで読む（改善版）
	decoder := json.NewDecoder(resp.Body)
	lastStatus := ""

	for {
		var progress struct {
			Status    string `json:"status"`
			Completed int64  `json:"completed,omitempty"`
			Total     int64  `json:"total,omitempty"`
		}

		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read pull progress: %w", err)
		}

		// 同じメッセージの連続を避ける
		currentStatus := progress.Status
		if progress.Total > 0 {
			percentage := float64(progress.Completed) / float64(progress.Total) * 100
			currentStatus = fmt.Sprintf("%s (%.1f%%)", progress.Status, percentage)
		}

		if currentStatus != lastStatus {
			log.Printf("📥 [Ollama] %s", currentStatus)
			lastStatus = currentStatus
		}
	}

	log.Printf("✅ [Ollama] Model '%s' pulled successfully", model)
	return nil
}

// waitForReady はOllamaが起動するまで待機します
func (c *ollamaClient) waitForReady(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	retryInterval := 2 * time.Second

	for time.Now().Before(deadline) {
		// コンテキストがキャンセルされていないか確認
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Health check エンドポイントを叩く
		url := fmt.Sprintf("%s/api/tags", c.baseURL)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create health check request: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			log.Println("✅ [Ollama] Service is ready")
			return nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		log.Printf("⏳ [Ollama] Waiting for service... (retrying in %v)", retryInterval)

		// リトライ前に待機（コンテキストキャンセルにも対応）
		select {
		case <-time.After(retryInterval):
			// 次のリトライへ
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout waiting for ollama (waited %v)", timeout)
}

// Generate は指定したモデルでテキスト生成を実行します
func (c *ollamaClient) Generate(ctx context.Context, model string, prompt string) (string, error) {
	// Step 1: リクエストボディ作成
	reqBody := OllamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // Non-streamingモード
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Step 2: HTTPリクエスト作成
	url := fmt.Sprintf("%s/api/generate", c.baseURL)
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Step 3: リクエスト送信
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Step 5: レスポンスデコード
	var response OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Step 6: 生成されたテキストを返す
	return response.Response, nil
}
