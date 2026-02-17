package client

type EmbedDocumentsRequest struct {
	Texts []string `json:"texts"`
}

type EmbedDocumentsResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
	Count      int         `json:"count"`
	Model      *string     `json:"model,omitempty"`
	Dim        *int        `json:"dim:omitempty"`
	ElapsedMs  *float64    `json:"elapsed_ms,omitempty"`
}

type GenerateRequest struct {
	Query     string   `json:"query"`
	Context   []string `json:"context"`
	MaxTokens int      `json:"max_tokens,omitempty"`
}

type GenerateResponse struct {
	Answer string `json:"answer"`
	Model  string `json:"model"`
}

type AIWorkerError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
