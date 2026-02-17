package service

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient はMinIOとの通信を担当
type MinIOClient struct {
	client *minio.Client
}

// NewMinIOClient は新しい MinIOClient を作成
func NewMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOClient{
		client: client,
	}, nil
}

// DownloadFile はMinIOからファイルをダウンロードします
func (m *MinIOClient) DownloadFile(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}

	// オブジェクトが存在するか確認
	_, err = object.Stat()
	if err != nil {
		object.Close()
		return nil, fmt.Errorf("object not found in MinIO: %w", err)
	}

	return object, nil
}

// UploadFile はMinIOにファイルをアップロードします（将来用）
func (m *MinIOClient) UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	return nil
}
