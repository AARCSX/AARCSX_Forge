package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
)

// ObjectRepositoryPostgres implements Repository using PostgreSQL.
type ObjectRepositoryPostgres struct {
	db *database.Postgres
}

// NewObjectRepositoryPostgres creates a new PostgreSQL object repository.
func NewObjectRepositoryPostgres(db *database.Postgres) Repository {
	return &ObjectRepositoryPostgres{db: db}
}

// SaveObjectMeta saves object metadata to the database.
func (r *ObjectRepositoryPostgres) SaveObjectMeta(ctx context.Context, obj ObjectMeta) (ObjectMeta, error) {
	query := `
		INSERT INTO object_metadata (id, tenant_id, object_key, file_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			object_key = EXCLUDED.object_key,
			file_name = EXCLUDED.file_name,
			updated_at = EXCLUDED.updated_at
		RETURNING id, tenant_id, object_key, file_name, created_at, updated_at
	`
	now := time.Now()
	if obj.ID == "" {
		obj.ID = uuid.New().String()
	}
	var createdAt, updatedAt time.Time
	err := r.db.DB.QueryRowContext(ctx, query, obj.ID, obj.TenantID, obj.ObjectKey, obj.FileName, now, now).
		Scan(&obj.ID, &obj.TenantID, &obj.ObjectKey, &obj.FileName, &createdAt, &updatedAt)
	if err != nil {
		return ObjectMeta{}, err
	}
	obj.CreatedAt = createdAt
	obj.UpdatedAt = updatedAt
	return obj, nil
}

// GetObjectMeta retrieves object metadata by ID and tenant ID.
func (r *ObjectRepositoryPostgres) GetObjectMeta(ctx context.Context, tenantID, objectID string) (ObjectMeta, error) {
	query := `
		SELECT id, tenant_id, object_key, file_name, created_at, updated_at
		FROM object_metadata
		WHERE tenant_id = $1 AND id = $2
	`
	var obj ObjectMeta
	var createdAt, updatedAt time.Time
	err := r.db.DB.QueryRowContext(ctx, query, tenantID, objectID).
		Scan(&obj.ID, &obj.TenantID, &obj.ObjectKey, &obj.FileName, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ObjectMeta{}, ErrObjectNotFound
		}
		return ObjectMeta{}, err
	}
	obj.CreatedAt = createdAt
	obj.UpdatedAt = updatedAt
	return obj, nil
}