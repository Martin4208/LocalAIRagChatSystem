package client

type CreateCollectionRequest struct {
	Vectors VectorConfig `json:"vectors"`
}

type VectorConfig struct {
	Size     int    `json:"size"`
	Distance string `json:"distance"`
}

type UpsertPointsRequest struct {
	Points []Point `json:"points"`
}

type Point struct {
	ID      string                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type SearchRequest struct {
	Vector      []float64 `json:"vector"`
	Limit       int       `json:"limit"`
	WithPayload bool      `json:"with_payload"`
}

type SearchResponse struct {
	Result []SearchResult `json:"result"`
}

type SearchResult struct {
	ID      string                 `json:"id"`
	Score   float64                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}
