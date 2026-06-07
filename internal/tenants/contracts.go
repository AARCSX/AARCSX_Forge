package tenants

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Repository defines the contract for tenant persistence.
type Repository interface {
	Create(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	Update(ctx context.Context, t *Tenant) error
}

// Service defines the contract for tenant business logic.
type Service interface {
	CreateTenant(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*TenantResponse, error)
	GetTenantBySlug(ctx context.Context, slug string) (*TenantResponse, error)
	UpdateTenant(ctx context.Context, id uuid.UUID, req *UpdateTenantRequest) (*TenantResponse, error)
}

// Errors
var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantExists   = errors.New("tenant already exists")
)