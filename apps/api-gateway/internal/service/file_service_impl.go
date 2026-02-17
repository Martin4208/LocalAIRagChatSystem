package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/storage"
	"github.com/google/uuid"
)

// FileServiceImpl implements FileService interface
type FileServiceImpl struct {
	queries       *db.Queries
	storageClient storage.ObjectStorageClient
	storageBucket string // MinIO bucket name for files
	qdrantClient  client.QdrantClient
}

// NewFileService creates a new FileService instance
func NewFileService(
	queries *db.Queries,
	storageClient storage.ObjectStorageClient,
	storageBucket string,
	qdrantClient client.QdrantClient,
) FileService {
	return &FileServiceImpl{
		queries:       queries,
		storageClient: storageClient,
		storageBucket: storageBucket,
	}
}

// UploadFile implements FileService.UploadFile
func (fs *FileServiceImpl) UploadFile(
	ctx context.Context,
	workspaceID uuid.UUID,
	file multipart.File,
	header *multipart.FileHeader,
	directoryID *uuid.UUID,
	tags []string,
) (*FileUploadResponse, error) {
	// Step 1: Calculate SHA256 hash while reading file
	hasher := sha256.New()
	fileContent := io.TeeReader(file, hasher)
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	sha256Hash := fmt.Sprintf("%x", hasher.Sum(nil))

	log.Printf("📁 File upload started: name=%s, hash=%s", header.Filename, sha256Hash)

	// Step 2: Check if file with same hash already exists
	existingFile, err := fs.queries.GetFileByHash(ctx, sha256Hash)
	if err == nil {
		log.Printf("🔄 File with same hash exists: file_id=%s, bucket=%s, key=%s",
			existingFile.ID, existingFile.MinioBucket, existingFile.MinioKey)

		// File already exists in MinIO, just create a new document reference
		return fs.createDocumentReference(
			ctx,
			workspaceID,
			existingFile.ID,
			header.Filename,
			directoryID,
			tags,
		)
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing file: %w", err)
	}

	// Step 3: Upload to MinIO
	minioKey := fs.generateMinIOKey(workspaceID, header.Filename)
	log.Printf("📤 Uploading to MinIO: bucket=%s, key=%s", fs.storageBucket, minioKey)

	err = fs.storageClient.PutObject(
		ctx,
		fs.storageBucket,
		minioKey,
		bytes.NewReader(fileBytes),
		int64(len(fileBytes)),
		storage.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	log.Printf("✅ MinIO upload successful: key=%s", minioKey)

	// Step 4: Create file record in database
	newFile, err := fs.queries.CreateFile(ctx, db.CreateFileParams{
		Sha256Hash: sha256Hash,
		MimeType:   header.Header.Get("Content-Type"),
		SizeBytes:  int64(len(fileBytes)),
		OriginalFilename: sql.NullString{
			String: header.Filename,
			Valid:  true,
		},
		MinioBucket: fs.storageBucket,
		MinioKey:    minioKey,
	})
	if err != nil {
		// Cleanup: remove from MinIO if DB insert fails
		_ = fs.storageClient.RemoveObject(ctx, fs.storageBucket, minioKey)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	log.Printf("✅ File record created: file_id=%s", newFile.ID)

	// Step 5: Create document reference
	return fs.createDocumentReference(
		ctx,
		workspaceID,
		newFile.ID,
		header.Filename,
		directoryID,
		tags,
	)
}

// ListFiles implements FileService.ListFiles
func (fs *FileServiceImpl) ListFiles(
	ctx context.Context,
	workspaceID uuid.UUID,
	filter FileListFilter,
) (*FileListResponse, error) {
	// Validate workspace exists
	_, err := fs.queries.GetWorkspace(ctx, workspaceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("workspace not found")
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	// Default limit and offset
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// Convert directoryID to uuid.NullUUID
	dirID := uuid.NullUUID{}
	if filter.DirectoryID != nil {
		dirID = uuid.NullUUID{
			UUID:  *filter.DirectoryID,
			Valid: true,
		}
	}

	// Query documents
	documents, err := fs.queries.ListDocuments(
		ctx,
		db.ListDocumentsParams{
			WorkspaceID: workspaceID,
			Offset:      int32(filter.Offset),
			Limit:       int32(filter.Limit),
			DirectoryID: dirID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// Count total
	total, err := fs.queries.CountDocuments(ctx, db.CountDocumentsParams{
		WorkspaceID: workspaceID,
		DirectoryID: dirID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// Map to response
	var files []FileMetadataResponse
	for _, doc := range documents {
		files = append(files, FileMetadataResponse{
			ID:        doc.ID,
			FileName:  doc.Name,
			MimeType:  doc.MimeType,
			SizeBytes: doc.SizeBytes,
			CreatedAt: doc.CreatedAt.Format(time.RFC3339),
			Status:    doc.Status,
		})
	}

	return &FileListResponse{
		Files: files,
		Total: total,
	}, nil
}

// GetFile implements FileService.GetFile
func (fs *FileServiceImpl) GetFile(
	ctx context.Context,
	workspaceID uuid.UUID,
	fileID uuid.UUID,
) (*FileMetadataResponse, error) {
	doc, err := fs.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:          fileID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("file not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return &FileMetadataResponse{
		ID:        doc.ID,
		FileName:  doc.Name,
		MimeType:  doc.MimeType,
		SizeBytes: doc.SizeBytes,
		CreatedAt: doc.CreatedAt.Format(time.RFC3339),
		Status:    doc.Status,
	}, nil
}

// DeleteFile implements FileService.DeleteFile (hard delete)
func (fs *FileServiceImpl) DeleteFile(
	ctx context.Context,
	workspaceID uuid.UUID,
	fileID uuid.UUID,
) error {
	// Step 1: 削除するドキュメント情報を取得
	documents, err := fs.queries.GetDocumentsByFileID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get documents: %w", err)
	}

	// Step 2: 各ドキュメントのベクトルをQdrantから削除
	collectionName := fmt.Sprintf("workspace_%s", workspaceID.String())
	for _, doc := range documents {
		err := fs.qdrantClient.DeletePointsByDocumentID(
			ctx,
			collectionName,
			doc.ID.String(),
		)
		if err != nil {
			// Qdrant削除失敗はログだけ（ベクトルが既に無い可能性もある）
			fmt.Printf("⚠️  Failed to delete from Qdrant: %v\n", err)
		}
	}

	// Step 3: documentsテーブルから削除
	err = fs.queries.HardDeleteDocumentsByFileID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}

	// Step 4: document_chunksも削除（CASCADE設定があれば不要だが明示的に）
	for _, doc := range documents {
		_ = fs.queries.DeleteDocumentChunks(ctx, doc.ID)
	}

	// Step 5: fileレコードを削除（MinIOは保持。複数documentで共有されている可能性）
	// 注: file削除はコメントアウト（deduplication対応のため）
	// err = fs.queries.DeleteFile(ctx, fileID)

	return nil
}

// DownloadFile implements FileService.DownloadFile
func (fs *FileServiceImpl) DownloadFile(
	ctx context.Context,
	workspaceID uuid.UUID,
	fileID uuid.UUID,
) (io.ReadCloser, error) {
	// Get document to find file location
	doc, err := fs.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:          fileID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("file not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Retrieve from MinIO
	reader, err := fs.storageClient.GetObject(
		ctx,
		doc.MinioBucket,
		doc.MinioKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download from MinIO: %w", err)
	}

	// Update last accessed timestamp
	_ = fs.queries.UpdateFileLastAccessed(ctx, doc.FileID)

	return reader, nil
}

// --- Helper Methods ---

// createDocumentReference creates a document entry for a file
func (fs *FileServiceImpl) createDocumentReference(
	ctx context.Context,
	workspaceID uuid.UUID,
	fileID uuid.UUID,
	fileName string,
	directoryID *uuid.UUID,
	tags []string,
) (*FileUploadResponse, error) {
	// Convert directoryID to uuid.NullUUID
	dirID := uuid.NullUUID{}
	if directoryID != nil {
		dirID = uuid.NullUUID{
			UUID:  *directoryID,
			Valid: true,
		}
	}

	doc, err := fs.queries.CreateDocument(ctx, db.CreateDocumentParams{
		WorkspaceID: workspaceID,
		DirectoryID: dirID,
		FileID:      fileID,
		Name:        fileName,
		Tags:        tags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Fetch file details
	file, err := fs.queries.GetFileById(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return &FileUploadResponse{
		ID:         doc.ID,
		FileName:   doc.Name,
		MimeType:   file.MimeType,
		SizeBytes:  file.SizeBytes,
		SHA256Hash: file.Sha256Hash,
		CreatedAt:  doc.CreatedAt.Format(time.RFC3339),
	}, nil
}

// generateMinIOKey creates a unique key path for MinIO storage
// Format: workspaceID/timestamp/uuid+ext
func (fs *FileServiceImpl) generateMinIOKey(workspaceID uuid.UUID, filename string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	return fmt.Sprintf("%s/%d/%s%s", workspaceID, timestamp, uuid.New(), ext)
}
