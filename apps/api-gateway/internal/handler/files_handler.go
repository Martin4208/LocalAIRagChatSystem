package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/service"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ListFiles implements GET /workspaces/{workspaceId}/files
func (h *Handler) ListFiles(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	params api.ListFilesParams,
) {
	// Parse parameters
	limit := int(50)
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	offset := int(0)
	if params.Offset != nil {
		offset = int(*params.Offset)
	}

	var directoryID *uuid.UUID
	if params.DirectoryId != nil {
		directoryID = (*uuid.UUID)(params.DirectoryId)
	}

	// Call service
	result, err := h.fileService.ListFiles(r.Context(), uuid.UUID(workspaceId), service.FileListFilter{
		DirectoryID: directoryID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list files", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UploadFile implements POST /workspaces/{workspaceId}/files/upload
func (h *Handler) UploadFile(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
) {
	// Parse multipart form (100MB max)
	err := r.ParseMultipartForm(100 * 1024 * 1024)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form", err)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field required", err)
		return
	}
	defer file.Close()

	// Parse optional directoryId
	var directoryID *uuid.UUID
	if dirIDStr := r.FormValue("directoryId"); dirIDStr != "" {
		if parsed, err := uuid.Parse(dirIDStr); err == nil {
			directoryID = &parsed
		}
	}

	// Parse optional tags
	var tags []string
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		json.Unmarshal([]byte(tagsStr), &tags)
	}

	// Call service
	result, err := h.fileService.UploadFile(
		r.Context(),
		uuid.UUID(workspaceId),
		file,
		header,
		directoryID,
		tags,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "upload failed", err)
		return
	}

	go func() {
		processCtx := context.Background() // 別のcontextで非同期実行

		opts := service.ProcessOptions{
			ChunkSize:      500,
			ChunkOverlap:   50,
			ForceReprocess: false,
		}

		err := h.documentProcessor.ProcessDocument(
			processCtx,
			uuid.UUID(workspaceId),
			result.ID,
			opts,
		)
		if err != nil {
			log.Printf("Failed to process document %s: %v", result.ID, err)
		} else {
			log.Printf("Successfully processed document %s", result.ID)
		}
	}()

	createdAt, err := time.Parse(time.RFC3339, result.CreatedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid createdAt", err)
		return
	}

	sha := result.SHA256Hash

	resp := api.FileUploadResponse{
		Id:         result.ID,
		FileName:   result.FileName,
		MimeType:   result.MimeType,
		SizeBytes:  result.SizeBytes,
		Sha256Hash: &sha,
		CreatedAt:  createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetFile implements GET /workspaces/{workspaceId}/files/{fileId}
func (h *Handler) GetFile(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	fileId openapi_types.UUID,
) {
	result, err := h.fileService.GetFile(r.Context(), uuid.UUID(workspaceId), uuid.UUID(fileId))
	if err != nil {
		writeError(w, http.StatusNotFound, "file not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeleteFile implements DELETE /workspaces/{workspaceId}/files/{fileId}
func (h *Handler) DeleteFile(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	fileId openapi_types.UUID,
) {
	err := h.fileService.DeleteFile(r.Context(), uuid.UUID(workspaceId), uuid.UUID(fileId))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DownloadFile implements GET /workspaces/{workspaceId}/files/{fileId}/download
func (h *Handler) DownloadFile(
	w http.ResponseWriter,
	r *http.Request,
	workspaceId openapi_types.UUID,
	fileId openapi_types.UUID,
) {
	reader, err := h.fileService.DownloadFile(r.Context(), uuid.UUID(workspaceId), uuid.UUID(fileId))
	if err != nil {
		writeError(w, http.StatusNotFound, "file not found", err)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, reader)
}
