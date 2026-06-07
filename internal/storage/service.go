package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/storage/provider"
)

// StorageService implements the Service interface for storage operations.
type StorageService struct {
	provider Provider
	repo     Repository
	logger   *logger.Logger
	cfg      *config.Config
}

// NewStorageService creates a new storage service instance.
func NewStorageService(cfg *config.Config, logger *logger.Logger) (*StorageService, error) {
	var p Provider
	var err error

	// Initialize storage provider based on configuration
	switch cfg.Storage.Provider {
	case "minio":
		p, err = provider.NewMinIOProvider(
			cfg.Storage.MinIO.Endpoint,
			cfg.Storage.MinIO.AccessKey,
			cfg.Storage.MinIO.SecretKey,
			cfg.Storage.MinIO.Bucket,
			cfg.Storage.MinIO.UseSSL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize minio provider: %w", err)
		}
	case "s3":
		return nil, errors.New("S3 provider not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Storage.Provider)
	}

	return &StorageService{
		provider: p,
		logger:   logger,
		cfg:      cfg,
	}, nil
}

// WithRepository sets the repository for the storage service.
func (s *StorageService) WithRepository(repo Repository) *StorageService {
	s.repo = repo
	return s
}

// CreateUploadURL generates a presigned URL for uploading a file and saves metadata.
func (s *StorageService) CreateUploadURL(ctx context.Context, in CreateUploadURLInput) (SignedURL, error) {
	// Validate input
	if in.TenantID == "" {
		return SignedURL{}, errors.New("tenant ID is required")
	}
	if in.ObjectID == "" {
		return SignedURL{}, errors.New("object ID is required")
	}
	if in.FileName == "" {
		return SignedURL{}, errors.New("file name is required")
	}
	if in.TTLSecond <= 0 {
		in.TTLSecond = 3600 // Default to 1 hour
	}

	// Generate object key with tenant isolation
	objectKey := fmt.Sprintf("%s/%s", in.TenantID, in.ObjectID)

	// Generate presigned PUT URL
	uploadURL, err := s.provider.PresignPut(ctx, objectKey, in.TTLSecond)
	if err != nil {
		s.logger.Sugar().Errorw("failed to generate upload URL", "error", err, "tenant_id", in.TenantID, "object_id", in.ObjectID)
		return SignedURL{}, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	// Save object metadata
	if s.repo != nil {
		objMeta := ObjectMeta{
			ID:        in.ObjectID,
			TenantID:  in.TenantID,
			ObjectKey: objectKey,
			FileName:  in.FileName,
		}
		_, err := s.repo.SaveObjectMeta(ctx, objMeta)
		if err != nil {
			s.logger.Sugar().Warnw("failed to save object metadata", "error", err, "tenant_id", in.TenantID, "object_id", in.ObjectID)
			// Continue anyway as the URL is still valid
		}
	}

	return SignedURL{URL: uploadURL}, nil
}

// CreateDownloadURL generates a presigned URL for downloading a file.
func (s *StorageService) CreateDownloadURL(ctx context.Context, in CreateDownloadURLInput) (SignedURL, error) {
	// Validate input
	if in.TenantID == "" {
		return SignedURL{}, errors.New("tenant ID is required")
	}
	if in.ObjectID == "" {
		return SignedURL{}, errors.New("object ID is required")
	}
	if in.TTLSecond <= 0 {
		in.TTLSecond = 3600 // Default to 1 hour
	}

	// Generate object key with tenant isolation
	objectKey := fmt.Sprintf("%s/%s", in.TenantID, in.ObjectID)

	// Generate presigned GET URL
	downloadURL, err := s.provider.PresignGet(ctx, objectKey, in.TTLSecond)
	if err != nil {
		s.logger.Sugar().Errorw("failed to generate download URL", "error", err, "tenant_id", in.TenantID, "object_id", in.ObjectID)
		return SignedURL{}, fmt.Errorf("failed to generate download URL: %w", err)
	}

	return SignedURL{URL: downloadURL}, nil
}