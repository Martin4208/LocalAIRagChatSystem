package handler

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// 1. Content-Typeヘッダーを設定
	w.Header().Set("Content-Type", "application/json")

	// 2. ステータスコードを設定（200 OK）
	w.WriteHeader(http.StatusOK)

	// 3. レスポンスボディを作成
	response := map[string]string{
		"status": "ok",
	}

	// 4. JSONに変換してwに書き込む
	json.NewEncoder(w).Encode(response)
}
