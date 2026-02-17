package storage

import (
	"context"
	"io"
)

// ObjectStorageClient defines the interface for object storage operations
// This abstraction allows us to swap MinIO with other S3-compatible services
type ObjectStorageClient interface {
	// PutObject uploads a file to storage
	// bucket: S3 bucket name
	// objectName: path/filename in the bucket
	// reader: file content
	// size: total bytes to upload (-1 if unknown)
	// returns: error if upload fails
	PutObject(
		ctx context.Context,
		bucket string,
		objectName string,
		reader io.Reader,
		size int64,
		opts PutObjectOptions,
	) error

	// GetObject downloads a file from storage
	// returns: io.ReadCloser that must be closed by caller
	GetObject(
		ctx context.Context,
		bucket string,
		objectName string,
	) (io.ReadCloser, error)

	// RemoveObject deletes a file from storage
	RemoveObject(
		ctx context.Context,
		bucket string,
		objectName string,
	) error

	// BucketExists checks if a bucket exists
	BucketExists(ctx context.Context, bucket string) (bool, error)

	// MakeBucket creates a new bucket
	MakeBucket(ctx context.Context, bucket string) error
}

// PutObjectOptions contains optional parameters for PutObject
type PutObjectOptions struct {
	ContentType string
	Metadata    map[string]string
}
