package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOProvider implements the Provider interface for MinIO object storage.
type MinIOProvider struct {
	client *minio.Client
	bucket string
	useSSL bool
}

// NewMinIOProvider creates a new MinIO provider instance.
func NewMinIOProvider(endpoint, accessKey, secretKey string, bucket string, useSSL bool) (*MinIOProvider, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket %s: %w", bucket, err)
		}
	}

	return &MinIOProvider{
		client: minioClient,
		bucket: bucket,
		useSSL: useSSL,
	}, nil
}

// PresignPut generates a presigned URL for uploading an object.
func (p *MinIOProvider) PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
	// Validate inputs
	if objectKey == "" {
		return "", fmt.Errorf("object key cannot be empty")
	}
	if ttlSeconds <= 0 {
		return "", fmt.Errorf("ttl seconds must be positive")
	}

	// Generate presigned PUT URL
	url, err := p.client.PresignedPutObject(ctx, p.bucket, objectKey, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL: %w", err)
	}

	return url.String(), nil
}

// PresignGet generates a presigned URL for downloading an object.
func (p *MinIOProvider) PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error) {
	// Validate inputs
	if objectKey == "" {
		return "", fmt.Errorf("object key cannot be empty")
	}
	if ttlSeconds <= 0 {
		return "", fmt.Errorf("ttl seconds must be positive")
	}

	// Generate presigned GET URL
	url, err := p.client.PresignedGetObject(ctx, p.bucket, objectKey, time.Duration(ttlSeconds)*time.Second, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned GET URL: %w", err)
	}

	return url.String(), nil
}