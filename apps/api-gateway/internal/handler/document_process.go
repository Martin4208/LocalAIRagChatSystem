package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/service"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ProcessDocument は、ドキュメントの処理を開始します
func (h *Handler) ProcessDocument(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	documentId openapi_types.UUID,
) {
	ctx := r.Context()

	// Step 1: リクエストボディを読み込む（オプション）
	var opts struct {
		ChunkSize      *int  `json:"chunk_size"`
		ChunkOverlap   *int  `json:"chunk_overlap"`
		ForceReprocess *bool `json:"force_reprocess"`
	}

	// リクエストボディがあれば読み込む
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}
	}

	// デフォルト値を設定
	chunkSize := 500
	if opts.ChunkSize != nil {
		chunkSize = *opts.ChunkSize
	}

	chunkOverlap := 50
	if opts.ChunkOverlap != nil {
		chunkOverlap = *opts.ChunkOverlap
	}

	forceReprocess := false
	if opts.ForceReprocess != nil {
		forceReprocess = *opts.ForceReprocess
	}

	// Step 2: Service層を呼ぶ
	processOpts := service.ProcessOptions{
		ChunkSize:      chunkSize,
		ChunkOverlap:   chunkOverlap,
		ForceReprocess: forceReprocess,
	}

	err := h.documentProcessor.ProcessDocument(ctx, workspaceId, documentId, processOpts)
	if err != nil {
		log.Printf("Failed to process document: %v", err)

		switch err {
		case service.ErrDocumentNotFound:
			respondError(w, http.StatusNotFound, "DOCUMENT_NOT_FOUND", "Document not found")
		case service.ErrAlreadyProcessing:
			respondError(w, http.StatusConflict, "ALREADY_PROCESSING", "Document is already being processed")
		default:
			respondError(w, http.StatusInternalServerError, "PROCESSING_ERROR", "Failed to process document")
		}
		return
	}

	// Step 3: レスポンスを返す（202 Accepted）
	respondJSON(w, http.StatusAccepted, map[string]string{
		"status":  "processing",
		"message": "Document processing started successfully",
	})
}

// GetDocumentStatus は、ドキュメントの処理状況を取得します
func (h *Handler) GetDocumentStatus(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	documentId openapi_types.UUID,
) {
	ctx := r.Context()

	status, err := h.documentProcessor.GetDocumentStatus(ctx, workspaceId, documentId)
	if err != nil {
		log.Printf("Failed to get document status: %v", err)

		if err == service.ErrDocumentNotFound {
			respondError(w, http.StatusNotFound, "DOCUMENT_NOT_FOUND", "Document not found")
		} else {
			respondError(w, http.StatusInternalServerError, "STATUS_ERROR", "Failed to get document status")
		}
		return
	}

	respondJSON(w, http.StatusOK, status)
}

// GetDocumentChunks は、ドキュメントのチャンク一覧を取得します
func (h *Handler) GetDocumentChunks(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	documentId openapi_types.UUID,
	params api.GetDocumentChunksParams, // ← パラメータ追加
) {
	ctx := r.Context()

	// パラメータからpage/limitを取得（デフォルト値付き）
	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}

	// Service層を呼ぶ
	chunks, total, err := h.documentProcessor.GetDocumentChunks(ctx, workspaceId, documentId, page, limit)
	if err != nil {
		log.Printf("Failed to get document chunks: %v", err)

		if err == service.ErrDocumentNotFound {
			respondError(w, http.StatusNotFound, "DOCUMENT_NOT_FOUND", "Document not found")
		} else {
			respondError(w, http.StatusInternalServerError, "CHUNKS_ERROR", "Failed to get document chunks")
		}
		return
	}

	// レスポンスを返す
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"chunks": chunks,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}
