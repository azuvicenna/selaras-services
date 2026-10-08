package client

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
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

func (s *MinioStorage) UploadFile(ctx context.Context, patientID, resourceID, fileName, contentType string, content io.Reader, size int64) (string, error) {
	key := path.Join("emr", patientID, resourceID, path.Base(fileName))

	_, err := s.client.PutObject(ctx, s.bucket, key, content, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *MinioStorage) GetPresignedURL(ctx context.Context, fileURL string, expiry time.Duration) (string, error) {
	if expiry <= 0 {
		expiry = 15 * time.Minute
	}
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, fileURL, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return presignedURL.String(), nil
}

func (s *MinioStorage) DeleteFile(ctx context.Context, fileURL string) error {
	return s.client.RemoveObject(ctx, s.bucket, fileURL, minio.RemoveObjectOptions{})
}
