package storage

import (
	"context"
	"errors"
	"time"
)

// Service defines the contract for storage business logic.
type Service interface {
	CreateUploadURL(ctx context.Context, in CreateUploadURLInput) (SignedURL, error)
	CreateDownloadURL(ctx context.Context, in CreateDownloadURLInput) (SignedURL, error)
}

// Provider defines the contract for storage providers.
type Provider interface {
	PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error)
	PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error)
}

// Repository defines the contract for storage persistence.
type Repository interface {
	SaveObjectMeta(ctx context.Context, obj ObjectMeta) (ObjectMeta, error)
	GetObjectMeta(ctx context.Context, tenantID, objectID string) (ObjectMeta, error)
}

// CreateUploadURLInput represents the input for creating an upload URL.
type CreateUploadURLInput struct {
	TenantID  string
	ObjectID  string
	FileName  string
	TTLSecond int
}

// CreateDownloadURLInput represents the input for creating a download URL.
type CreateDownloadURLInput struct {
	TenantID  string
	ObjectID  string
	TTLSecond int
}

// ObjectMeta represents object metadata.
type ObjectMeta struct {
	ID        string
	TenantID  string
	ObjectKey string
	FileName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Errors
var (
	ErrObjectNotFound = errors.New("object not found")
)

// SignedURL represents a signed URL for upload/download operations.
type SignedURL struct {
	URL string
}
