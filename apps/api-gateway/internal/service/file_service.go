package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
)

// FileMetadataResponse represents file metadata returned to clients
type FileMetadataResponse struct {
	ID         uuid.UUID
	FileName   string
	MimeType   string
	SizeBytes  int64
	SHA256Hash string
	CreatedAt  string // RFC3339 format
	Tags       []string
	Status     string
}

// FileUploadResponse represents the response after file upload
type FileUploadResponse struct {
	ID         uuid.UUID
	FileName   string
	MimeType   string
	SizeBytes  int64
	SHA256Hash string
	CreatedAt  string
}

// FileListResponse represents the response for file listing
type FileListResponse struct {
	Files []FileMetadataResponse
	Total int64
}

// FileListFilter represents filter options for file listing
type FileListFilter struct {
	DirectoryID *uuid.UUID
	Limit       int
	Offset      int
}

// FileService defines the business logic for file operations
type FileService interface {
	// UploadFile handles file upload with deduplication by SHA256
	// Returns FileUploadResponse on success
	UploadFile(
		ctx context.Context,
		workspaceID uuid.UUID,
		file multipart.File,
		header *multipart.FileHeader,
		directoryID *uuid.UUID,
		tags []string,
	) (*FileUploadResponse, error)

	// ListFiles retrieves files in a workspace with optional filtering
	ListFiles(
		ctx context.Context,
		workspaceID uuid.UUID,
		filter FileListFilter,
	) (*FileListResponse, error)

	// GetFile retrieves metadata for a single file
	GetFile(
		ctx context.Context,
		workspaceID uuid.UUID,
		fileID uuid.UUID,
	) (*FileMetadataResponse, error)

	// DeleteFile performs soft delete on a file
	DeleteFile(
		ctx context.Context,
		workspaceID uuid.UUID,
		fileID uuid.UUID,
	) error

	// DownloadFile retrieves file content from MinIO
	// Caller is responsible for closing the io.ReadCloser
	DownloadFile(
		ctx context.Context,
		workspaceID uuid.UUID,
		fileID uuid.UUID,
	) (io.ReadCloser, error)
}
