package notifications

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
)

// DeliveryJobPostgres implements Repository using PostgreSQL.
type DeliveryJobPostgres struct {
	db *database.Postgres
}

// NewDeliveryJobPostgres creates a new PostgreSQL delivery job repository.
func NewDeliveryJobPostgres(db *database.Postgres) Repository {
	return &DeliveryJobPostgres{db: db}
}

// CreateDelivery creates a new delivery job record.
func (r *DeliveryJobPostgres) CreateDelivery(ctx context.Context, in DeliveryJob) (DeliveryJob, error) {
	query := `
		INSERT INTO notification_deliveries (id, tenant_id, channel, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, tenant_id, channel, status, created_at, updated_at
	`
	now := time.Now()
	if in.ID == "" {
		in.ID = uuid.New().String()
	}
	var createdAt, updatedAt time.Time
	err := r.db.DB.QueryRowContext(ctx, query, in.ID, in.TenantID, in.Channel, in.Status, now, now).
		Scan(&in.ID, &in.TenantID, &in.Channel, &in.Status, &createdAt, &updatedAt)
	if err != nil {
		return DeliveryJob{}, err
	}
	in.CreatedAt = createdAt
	in.UpdatedAt = updatedAt
	return in, nil
}

// MarkDelivered updates the delivery job status to delivered.
func (r *DeliveryJobPostgres) MarkDelivered(ctx context.Context, jobID string) error {
	query := `
		UPDATE notification_deliveries
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	now := time.Now()
	_, err := r.db.DB.ExecContext(ctx, query, "delivered", now, jobID)
	return err
}

// MarkFailed updates the delivery job status to failed with a reason.
func (r *DeliveryJobPostgres) MarkFailed(ctx context.Context, jobID string, reason string) error {
	query := `
		UPDATE notification_deliveries
		SET status = $1, failure_reason = $2, updated_at = $3
		WHERE id = $4
	`
	now := time.Now()
	_, err := r.db.DB.ExecContext(ctx, query, "failed", reason, now, jobID)
	return err
}