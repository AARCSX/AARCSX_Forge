package tenants

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
)

// TenantService provides tenant business logic.
type TenantService struct {
	repo   Repository
	logger *logger.Logger
}

// NewTenantService creates a new tenant service.
func NewTenantService(repo Repository, logger *logger.Logger) *TenantService {
	return &TenantService{repo: repo, logger: logger}
}

// CreateTenant creates a new tenant.
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error) {
	// TODO: Add validation (e.g., slug uniqueness) via repo or additional logic.
	now := time.Now()
	t := &Tenant{
		ID:        uuid.New(),
		Name:      req.Name,
		Slug:      req.Slug,
		Status:    req.Status,
		Settings:  []byte("{}"), // default empty JSON
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return &TenantResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Slug:      t.Slug,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetTenantByID returns a tenant by ID.
func (s *TenantService) GetTenantByID(ctx context.Context, id uuid.UUID) (*TenantResponse, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &TenantResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Slug:      t.Slug,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetTenantBySlug returns a tenant by slug.
func (s *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*TenantResponse, error) {
	t, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return &TenantResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Slug:      t.Slug,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdateTenant updates a tenant.
func (s *TenantService) UpdateTenant(ctx context.Context, id uuid.UUID, req *UpdateTenantRequest) (*TenantResponse, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Slug != "" {
		existing.Slug = req.Slug
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return &TenantResponse{
		ID:        existing.ID.String(),
		Name:      existing.Name,
		Slug:      existing.Slug,
		Status:    existing.Status,
		CreatedAt: existing.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: existing.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}