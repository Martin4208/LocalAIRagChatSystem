package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient wraps minio-go client and implements ObjectStorageClient
type MinIOClient struct {
	client *minio.Client
}

// NewMinIOClient creates a new MinIO client
// endpoint: MinIO server address (e.g., "localhost:9000")
// accessKey: MinIO access key
// secretKey: MinIO secret key
// useSSL: whether to use HTTPS
func NewMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{
		client: client,
	}, nil
}

// PutObject implements ObjectStorageClient.PutObject
func (m *MinIOClient) PutObject(
	ctx context.Context,
	bucket string,
	objectName string,
	reader io.Reader,
	size int64,
	opts PutObjectOptions,
) error {
	_, err := m.client.PutObject(
		ctx,
		bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType:  opts.ContentType,
			UserMetadata: opts.Metadata,
		},
	)
	return err
}

// GetObject implements ObjectStorageClient.GetObject
func (m *MinIOClient) GetObject(
	ctx context.Context,
	bucket string,
	objectName string,
) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
}

// RemoveObject implements ObjectStorageClient.RemoveObject
func (m *MinIOClient) RemoveObject(
	ctx context.Context,
	bucket string,
	objectName string,
) error {
	return m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

// BucketExists implements ObjectStorageClient.BucketExists
func (m *MinIOClient) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return m.client.BucketExists(ctx, bucket)
}

// MakeBucket implements ObjectStorageClient.MakeBucket
func (m *MinIOClient) MakeBucket(ctx context.Context, bucket string) error {
	return m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}
