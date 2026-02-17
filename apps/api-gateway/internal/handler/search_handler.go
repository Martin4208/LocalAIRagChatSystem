// apps/api-gateway/internal/handler/search_handler.go

package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// SearchWorkspace は RAG 検索を実行
func (h *Handler) SearchWorkspace(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
) {
	ctx := r.Context()

	// リクエストボディを読み込む
	var reqBody api.SearchWorkspaceJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// バリデーション
	if reqBody.Query == "" {
		respondError(w, http.StatusBadRequest, "INVALID_QUERY", "Query cannot be empty")
		return
	}

	// デフォルト値設定
	topK := 5
	if reqBody.TopK != nil {
		topK = *reqBody.TopK
	}

	// Workspace の存在確認
	_, err := h.queries.GetWorkspace(ctx, workspaceId)
	if err != nil {
		respondError(w, http.StatusNotFound, "WORKSPACE_NOT_FOUND", "Workspace not found")
		return
	}

	// 検索実行
	log.Printf("Executing search: workspace=%s, query=%s", workspaceId, reqBody.Query)
	result, err := h.searchService.Search(ctx, workspaceId, reqBody.Query, topK)
	if err != nil {
		log.Printf("Search failed: %v", err)
		respondError(w, http.StatusInternalServerError, "SEARCH_ERROR", "Failed to execute search")
		return
	}

	// レスポンス作成
	response := map[string]interface{}{
		"answer":  result.Answer,
		"sources": result.Sources,
	}

	respondJSON(w, http.StatusOK, response)
}
