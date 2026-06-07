package observability

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AARCSX/AARCSX_Forge/internal/database"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"

)

// AuditLogRepository handles persistence of audit log entries.
// Repository layer only talks to DB - no business logic or request handling.
type AuditLogRepository struct {
	db     *database.Postgres
	logger *logger.Logger
}

// NewAuditLogRepository creates a new audit log repository.
func NewAuditLogRepository(db *database.Postgres, logger *logger.Logger) *AuditLogRepository {
	return &AuditLogRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new audit log entry.
func (r *AuditLogRepository) Create(ctx context.Context, log AuditLog) (int64, error) {
	query := `
		INSERT INTO audit_logs (
			tenant_id, user_id, action, resource, resource_id,
			details, ip_address, user_agent, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var id int64
	err := r.db.DB.QueryRowContext(
		ctx,
		query,
		log.TenantID,
		log.UserID,
		log.Action,
		log.Resource,
		log.ResourceID,
		log.Details,
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	).Scan(&id)

    if err != nil {
        r.logger.Error(ctx, "Failed to create audit log", err, map[string]any{})
        return 0, err
    }

    r.logger.Info(ctx, "Audit log created", map[string]any{"id": id})
    return id, nil
}

// GetByID retrieves an audit log entry by ID.
func (r *AuditLogRepository) GetByID(ctx context.Context, id int64) (*AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, resource, resource_id,
		       details, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE id = $1
	`

	var log AuditLog
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.TenantID,
		&log.UserID,
		&log.Action,
		&log.Resource,
		&log.ResourceID,
		&log.Details,
		&log.IPAddress,
		&log.UserAgent,
		&log.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error(ctx, "Failed to get audit log by ID", err, map[string]any{"id": id})
		return nil, err
	}

	return &log, nil
}

// List retrieves paginated audit log entries with optional filtering.
func (r *AuditLogRepository) List(ctx context.Context, tenantID int64, userID int64,
	action string, resource string, page int, pageSize int) ([]AuditLog, int64, error) {

	// Build WHERE clause
	where := []string{"1=1"}
	args := []interface{}{}
	argID := 1

	if tenantID != 0 {
		where = append(where, fmt.Sprintf("tenant_id = $%d", argID))
		args = append(args, tenantID)
		argID++
	}
	if userID != 0 {
		where = append(where, fmt.Sprintf("user_id = $%d", argID))
		args = append(args, userID)
		argID++
	}
	if action != "" {
		where = append(where, fmt.Sprintf("action = $%d", argID))
		args = append(args, action)
		argID++
	}
	if resource != "" {
		where = append(where, fmt.Sprintf("resource = $%d", argID))
		args = append(args, resource)
		argID++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs WHERE %s", strings.Join(where, " AND "))
	var totalCount int64
	err := r.db.DB.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		r.logger.Error(ctx, "Failed to count audit logs", err, map[string]any{})
		return nil, 0, err
	}

	// List query
	listQuery := fmt.Sprintf(`
		SELECT id, tenant_id, user_id, action, resource, resource_id,
		       details, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), argID, argID+1)

	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.DB.QueryContext(ctx, listQuery, args...)
	if err != nil {
		r.logger.Error(ctx, "Failed to list audit logs", err, map[string]any{})
		return nil, 0, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&log.Action,
			&log.Resource,
			&log.ResourceID,
			&log.Details,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			r.logger.Error(ctx, "Failed to scan audit log row", err, map[string]any{})
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(ctx, "Error iterating audit log rows", err, map[string]any{})
		return nil, 0, err
	}

	return logs, totalCount, nil
}