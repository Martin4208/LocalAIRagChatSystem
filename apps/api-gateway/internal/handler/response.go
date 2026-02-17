package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
)

// respondJSON はJSONレスポンスを返すヘルパー関数
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// respondError はOpenAPIのErrorスキーマに従ったエラーレスポンスを返す
func respondError(w http.ResponseWriter, status int, code, message string) {
	errorCode := code
	errorMessage := message
	timestamp := time.Now().Format(time.RFC3339)

	respondJSON(w, status, api.Error{
		Error: &struct {
			Code      *string                 `json:"code,omitempty"`
			Details   *map[string]interface{} `json:"details,omitempty"`
			Message   *string                 `json:"message,omitempty"`
			RequestId *string                 `json:"request_id,omitempty"`
			Timestamp *string                 `json:"timestamp,omitempty"`
		}{
			Code:      &errorCode,
			Message:   &errorMessage,
			Timestamp: &timestamp,
			Details:   nil,
			RequestId: nil,
		},
	})
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}

// intPtr は整数のポインタを返すヘルパー関数
func intPtr(i int) *int {
	return &i
}
