package tenants

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
)


// TenantRepositoryPostgres implements Repository using PostgreSQL.
type TenantRepositoryPostgres struct {
	db *database.Postgres
}

// NewTenantRepositoryPostgres creates a new PostgreSQL tenant repository.
func NewTenantRepositoryPostgres(db *database.Postgres) Repository {
	return &TenantRepositoryPostgres{db: db}
}

// Create inserts a new tenant into the database.
func (r *TenantRepositoryPostgres) Create(ctx context.Context, t *Tenant) error {
	query := `
		INSERT INTO tenants (id, name, slug, status, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	_, err := r.db.DB.ExecContext(ctx, query, t.ID, t.Name, t.Slug, t.Status, t.Settings, t.CreatedAt, t.UpdatedAt)
	return err
}

// GetByID retrieves a tenant by its ID.
func (r *TenantRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	query := `
		SELECT id, name, slug, status, settings, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`
	var t Tenant
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Slug, &t.Status, &t.Settings, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &t, nil
}

// GetBySlug retrieves a tenant by its slug.
func (r *TenantRepositoryPostgres) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	query := `
		SELECT id, name, slug, status, settings, created_at, updated_at
		FROM tenants
		WHERE slug = $1
	`
	var t Tenant
	err := r.db.DB.QueryRowContext(ctx, query, slug).Scan(
		&t.ID, &t.Name, &t.Slug, &t.Status, &t.Settings, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &t, nil
}

// Update updates an existing tenant.
func (r *TenantRepositoryPostgres) Update(ctx context.Context, t *Tenant) error {
	query := `
		UPDATE tenants
		SET name = $1, slug = $2, status = $3, settings = $4, updated_at = $5
		WHERE id = $6
	`
	t.UpdatedAt = time.Now()
	_, err := r.db.DB.ExecContext(ctx, query, t.Name, t.Slug, t.Status, t.Settings, t.UpdatedAt, t.ID)
	return err
}