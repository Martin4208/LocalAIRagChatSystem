package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sqlc-dev/pqtype"
)

// TODO: 実装する
func (h *Handler) ListWorkspaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspaces, err := h.queries.ListWorkspaces(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch workspaces",
		})
		return
	}

	// 変換
	apiWorkspaces := make([]api.Workspace, len(workspaces))
	for i, w := range workspaces {
		apiWorkspaces[i] = convertToAPIWorkspace(w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"workspaces": apiWorkspaces,
	}

	json.NewEncoder(w).Encode(response)
}

// TODO: 実装する
func (h *Handler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody api.CreateWorkspaceJSONBody

	// JSONをデコード
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// nameが空でないかチェック
	if reqBody.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Name is required",
		})
		return
	}

	// Descriptionの変換（*string → sql.NullString）
	var desc sql.NullString
	if reqBody.Description != nil {
		desc = sql.NullString{String: *reqBody.Description, Valid: true}
	} else {
		desc = sql.NullString{Valid: false}
	}

	params := db.CreateWorkspaceParams{
		Name:        reqBody.Name,
		Description: desc,
		Settings:    pqtype.NullRawMessage{Valid: false}, // 今は空
	}

	workspace, err := h.queries.CreateWorkspace(ctx, params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create workspace",
		})
		return
	}
	// レスポンス用に変換
	response := convertToAPIWorkspace(workspace)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// TODO: 実装する
func (h *Handler) GetWorkspace(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID) {
	ctx := r.Context()

	workspace, err := h.queries.GetWorkspace(ctx, workspaceId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Workspace not found",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch workspace",
		})
		return
	}

	response := convertToAPIWorkspace(workspace)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TODO: 実装する
func (h *Handler) UpdateWorkspace(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID) {
	ctx := r.Context()

	var reqBody api.UpdateWorkspaceJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// 既存データを取得
	existing, err := h.queries.GetWorkspace(ctx, workspaceId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Workspace not found",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch workspace",
		})
		return
	}

	// 更新するフィールドを決定
	name := existing.Name
	if reqBody.Name != nil {
		name = *reqBody.Name
	}

	var desc sql.NullString
	if reqBody.Description != nil {
		desc = sql.NullString{String: *reqBody.Description, Valid: true}
	} else {
		desc = existing.Description
	}

	params := db.UpdateWorkspaceParams{
		ID:          workspaceId,
		Name:        name,
		Description: desc,
	}

	err = h.queries.UpdateWorkspace(ctx, params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update workspace",
		})
		return
	}

	// 更新後のデータを取得
	updated, err := h.queries.GetWorkspace(ctx, workspaceId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch updated workspace",
		})
		return
	}

	response := convertToAPIWorkspace(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TODO: 実装する
func (h *Handler) DeleteWorkspace(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID) {
	ctx := r.Context()

	err := h.queries.DeleteWorkspace(ctx, workspaceId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to delete workspace",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func convertToAPIWorkspace(w db.Workspace) api.Workspace {
	var desc *string
	if w.Description.Valid {
		desc = &w.Description.String
	}

	var settings *map[string]interface{}
	if w.Settings.Valid && len(w.Settings.RawMessage) > 0 {
		var s map[string]interface{}
		json.Unmarshal(w.Settings.RawMessage, &s)
		settings = &s
	}

	return api.Workspace{
		Id:          w.ID,
		Name:        w.Name,
		Description: desc,
		Settings:    settings,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}
