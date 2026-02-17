package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ========================================
// ListAnalyses - GET /workspaces/{workspaceId}/analyses
// ========================================

func (h *Handler) ListAnalyses(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	params api.ListAnalysesParams,
) {
	ctx := r.Context()

	// パラメータのデフォルト値設定
	limit := 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	// ステータスフィルタ
	var status *string
	if params.Status != nil {
		statusStr := string(*params.Status)
		status = &statusStr
	}

	// Service層を呼び出し
	analyses, total, err := h.analysisService.ListAnalyses(
		ctx,
		workspaceId,
		status,
		limit,
		offset,
	)
	if err != nil {
		log.Printf("Failed to list analyses: %v", err)
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to list analyses")
		return
	}

	// レスポンス変換
	apiAnalyses := make([]api.Analysis, len(analyses))
	for i, analysis := range analyses {
		apiAnalyses[i] = listAnalysesRowToAPI(analysis)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"analyses": apiAnalyses,
		"total":    total,
	})
}

// ========================================
// CreateAnalysis - POST /workspaces/{workspaceId}/analyses
// ========================================

func (h *Handler) CreateAnalysis(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
) {
	ctx := r.Context()

	// リクエストボディをデコード
	var reqBody api.CreateAnalysisJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// バリデーション
	if reqBody.Title == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Title is required")
		return
	}

	if reqBody.AnalysisType == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Analysis type is required")
		return
	}

	// Config の変換（api.CreateAnalysis_Config → map[string]interface{}）
	var configMap map[string]interface{}
	if reqBody.Config != nil {
		// Config をマップに変換
		configBytes, _ := json.Marshal(reqBody.Config)
		json.Unmarshal(configBytes, &configMap)
	}

	// Service層を呼び出し
	analysis, err := h.analysisService.CreateAnalysis(
		ctx,
		workspaceId,
		reqBody.Title,
		reqBody.Description,
		string(reqBody.AnalysisType),
		configMap,
	)
	if err != nil {
		log.Printf("Failed to create analysis: %v", err)
		respondError(w, http.StatusInternalServerError, "CREATE_ERROR", "Failed to create analysis")
		return
	}

	// 202 Accepted で返す（バックグラウンド処理中）
	respondJSON(w, http.StatusAccepted, analysisToAPI(*analysis))
}

// ========================================
// GetAnalysis - GET /workspaces/{workspaceId}/analyses/{analysisId}
// ========================================

func (h *Handler) GetAnalysis(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	analysisId openapi_types.UUID,
) {
	ctx := r.Context()

	// Service層を呼び出し
	analysis, err := h.analysisService.GetAnalysis(ctx, workspaceId, analysisId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Analysis not found")
			return
		}
		log.Printf("Failed to get analysis: %v", err)
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to get analysis")
		return
	}

	respondJSON(w, http.StatusOK, analysisToAPI(*analysis))
}

// ========================================
// DeleteAnalysis - DELETE /workspaces/{workspaceId}/analyses/{analysisId}
// ========================================

func (h *Handler) DeleteAnalysis(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	analysisId openapi_types.UUID,
) {
	ctx := r.Context()

	// Service層を呼び出し
	err := h.analysisService.DeleteAnalysis(ctx, workspaceId, analysisId)
	if err != nil {
		log.Printf("Failed to delete analysis: %v", err)
		respondError(w, http.StatusInternalServerError, "DELETE_ERROR", "Failed to delete analysis")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ========================================
// GetAnalysisResults - GET /workspaces/{workspaceId}/analyses/{analysisId}/results
// ========================================

func (h *Handler) GetAnalysisResults(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	analysisId openapi_types.UUID,
) {
	ctx := r.Context()

	// Step 1: 分析が存在するか確認
	analysis, err := h.analysisService.GetAnalysis(ctx, workspaceId, analysisId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Analysis not found")
			return
		}
		log.Printf("Failed to get analysis: %v", err)
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to get analysis")
		return
	}

	// Step 2: ステータスチェック（completed 以外は 409 Conflict）
	if analysis.Status != "completed" {
		respondError(w, http.StatusConflict, "NOT_COMPLETED", "Analysis is not completed yet")
		return
	}

	// Step 3: 結果を取得
	results, err := h.analysisService.GetAnalysisResults(ctx, analysisId)
	if err != nil {
		log.Printf("Failed to get analysis results: %v", err)
		respondError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to get results")
		return
	}

	// レスポンス変換
	apiResults := make([]api.AnalysisResult, len(results))
	for i, result := range results {
		apiResults[i] = analysisResultToAPI(result)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"results": apiResults,
	})
}

// ========================================
// Helper Functions
// ========================================

// analysisToAPI は db.Analysis を api.Analysis に変換
func analysisToAPI(analysis db.Analysis) api.Analysis {
	var description *string
	if analysis.Description.Valid {
		description = &analysis.Description.String
	}

	// Config の変換
	var config *map[string]interface{}
	if analysis.Config.Valid && len(analysis.Config.RawMessage) > 0 {
		var configMap map[string]interface{}
		if err := json.Unmarshal(analysis.Config.RawMessage, &configMap); err == nil {
			config = &configMap
		}
	}

	// StartedAt の変換（time.Time 型）
	var startedAt *time.Time
	if analysis.StartedAt.Valid {
		startedAt = &analysis.StartedAt.Time
	}

	// CompletedAt の変換（time.Time 型）
	var completedAt *time.Time
	if analysis.CompletedAt.Valid {
		completedAt = &analysis.CompletedAt.Time
	}

	var errorMessage *string
	if analysis.ErrorMessage.Valid {
		errorMessage = &analysis.ErrorMessage.String
	}

	return api.Analysis{
		Id:           analysis.ID,
		WorkspaceId:  analysis.WorkspaceID,
		Title:        analysis.Title,
		Description:  description,
		AnalysisType: api.AnalysisAnalysisType(analysis.AnalysisType),
		Status:       api.AnalysisStatus(analysis.Status),
		Config:       config,
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
		ErrorMessage: errorMessage,
		CreatedAt:    analysis.CreatedAt,
		UpdatedAt:    analysis.UpdatedAt,
	}
}

// analysisResultToAPI は db.AnalysisResult を api.AnalysisResult に変換
func analysisResultToAPI(result db.AnalysisResult) api.AnalysisResult {
	// Content の変換（json.RawMessage → *map[string]interface{}）
	var content *map[string]interface{}
	if len(result.Content) > 0 {
		var contentMap map[string]interface{}
		if err := json.Unmarshal(result.Content, &contentMap); err == nil {
			content = &contentMap
		}
	}

	var imageUrl *string
	if result.ImageUrl.Valid {
		imageUrl = &result.ImageUrl.String
	}

	var metadata *map[string]interface{}
	if result.Metadata.Valid && len(result.Metadata.RawMessage) > 0 {
		var metaMap map[string]interface{}
		if err := json.Unmarshal(result.Metadata.RawMessage, &metaMap); err == nil {
			metadata = &metaMap
		}
	}

	return api.AnalysisResult{
		Id:         result.ID,
		AnalysisId: result.AnalysisID,
		ResultType: api.AnalysisResultResultType(result.ResultType),
		Content:    content,
		ImageUrl:   imageUrl,
		Metadata:   metadata,
		CreatedAt:  result.CreatedAt,
	}
}

// listAnalysesRowToAPI は db.ListAnalysesRow を api.Analysis に変換
func listAnalysesRowToAPI(row db.ListAnalysesRow) api.Analysis {
	var description *string
	if row.Description.Valid {
		description = &row.Description.String
	}

	// Config の変換
	var config *map[string]interface{}
	if row.Config.Valid && len(row.Config.RawMessage) > 0 {
		var configMap map[string]interface{}
		if err := json.Unmarshal(row.Config.RawMessage, &configMap); err == nil {
			config = &configMap
		}
	}

	// StartedAt の変換
	var startedAt *time.Time
	if row.StartedAt.Valid {
		startedAt = &row.StartedAt.Time
	}

	// CompletedAt の変換
	var completedAt *time.Time
	if row.CompletedAt.Valid {
		completedAt = &row.CompletedAt.Time
	}

	var errorMessage *string
	if row.ErrorMessage.Valid {
		errorMessage = &row.ErrorMessage.String
	}

	return api.Analysis{
		Id:           row.ID,
		WorkspaceId:  row.WorkspaceID,
		Title:        row.Title,
		Description:  description,
		AnalysisType: api.AnalysisAnalysisType(row.AnalysisType),
		Status:       api.AnalysisStatus(row.Status),
		Config:       config,
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
		ErrorMessage: errorMessage,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
