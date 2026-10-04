package client

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/azuvicenna/selaras-services/patient-service/internal/domain"
)

var _ domain.StorageClient = (*MinioStorage)(nil)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(ctx context.Context, endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioStorage, error) {
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioStorage{client: mc, bucket: bucket}, nil
}

// UploadFile menyimpan berkas dan mengembalikan object key, yang disimpan sebagai file_url.
func (s *MinioStorage) UploadFile(ctx context.Context, patientID, docID, fileName, contentType string, content io.Reader) (string, error) {
	key := path.Join("patients", patientID, docID, path.Base(fileName))

	_, err := s.client.PutObject(ctx, s.bucket, key, content, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return key, nil
}

// DeleteFile menerima object key yang dikembalikan UploadFile.
func (s *MinioStorage) DeleteFile(ctx context.Context, fileURL string) error {
	return s.client.RemoveObject(ctx, s.bucket, fileURL, minio.RemoveObjectOptions{})
}