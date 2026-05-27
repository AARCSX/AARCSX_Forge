package storage

import "context"

type Service interface {
	CreateUploadURL(ctx context.Context, in CreateUploadURLInput) (SignedURL, error)
	CreateDownloadURL(ctx context.Context, in CreateDownloadURLInput) (SignedURL, error)
}

type Provider interface {
	PresignPut(ctx context.Context, objectKey string, ttlSeconds int) (string, error)
	PresignGet(ctx context.Context, objectKey string, ttlSeconds int) (string, error)
}

type Repository interface {
	SaveObjectMeta(ctx context.Context, obj ObjectMeta) (ObjectMeta, error)
	GetObjectMeta(ctx context.Context, tenantID, objectID string) (ObjectMeta, error)
}

type CreateUploadURLInput struct {
	TenantID  string
	ObjectID  string
	FileName  string
	TTLSecond int
}

type CreateDownloadURLInput struct {
	TenantID  string
	ObjectID  string
	TTLSecond int
}

type ObjectMeta struct {
	ID        string
	TenantID  string
	ObjectKey string
	FileName  string
}

type SignedURL struct {
	URL string
}
