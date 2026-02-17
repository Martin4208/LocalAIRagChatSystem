package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AIWorkerClient interface {
	EmbedDocuments(ctx context.Context, texts []string) (*EmbedDocumentsResponse, error)
	EmbedQuery(ctx context.Context, text string) ([]float64, error)
	GenerateAnswer(ctx context.Context, query string, contextTexts []string) (*GenerateResponse, error)
}

type aiWorkerClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAIWorkerClient(baseURL string) AIWorkerClient {
	return &aiWorkerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *aiWorkerClient) EmbedDocuments(ctx context.Context, texts []string) (*EmbedDocumentsResponse, error) {
	reqBody := EmbedDocumentsRequest{
		Texts: texts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/api/v1/embed",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp AIWorkerError
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("AI Worker returned status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("AI Worker error: %s - %s", errResp.Code, errResp.Message)
	}

	var response EmbedDocumentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (c *aiWorkerClient) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	reqBody := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/api/v1/embed/query",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI Worker returned status %d", resp.StatusCode)
	}

	var response struct {
		Embedding []float64 `json:"embedding"`
		Dim       int       `json:"dim"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Embedding, nil
}

func (c *aiWorkerClient) GenerateAnswer(
	ctx context.Context,
	query string,
	contextTexts []string,
) (*GenerateResponse, error) {
	reqBody := GenerateRequest{
		Query:     query,
		Context:   contextTexts,
		MaxTokens: 500,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/api/v1/generate",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp AIWorkerError
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("AI Worker returned status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("AI Worker error: %s - %s", errResp.Code, errResp.Message)
	}

	var response GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
