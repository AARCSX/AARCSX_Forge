package identity

import (
	"context"

	"github.com/AARCSX/AARCSX_Forge/internal/database"
	"github.com/google/uuid"
)

// RBACPermissionChecker implements role-based access control.
type RBACPermissionChecker struct {
	db *database.Postgres
}

// NewRBACPermissionChecker creates a new RBAC permission checker.
func NewRBACPermissionChecker(db *database.Postgres) *RBACPermissionChecker {
	return &RBACPermissionChecker{db: db}
}

// Can checks if the actor has permission to perform the action on the tenant.
func (c *RBACPermissionChecker) Can(ctx context.Context, actorID, tenantID string, action string) (bool, error) {
	actorUUID, err := parseUUID(actorID)
	if err != nil {
		return false, err
	}
	tenantUUID, err := parseUUID(tenantID)
	if err != nil {
		return false, err
	}

	// Check if user has any role that grants the permission
	query := `
		SELECT COUNT(*) > 0
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1
		  AND r.tenant_id = $2
		  AND r.permissions ? $3
	`

	var hasPermission bool
	err = c.db.DB.QueryRowContext(ctx, query, actorUUID, tenantUUID, action).Scan(&hasPermission)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

// parseUUID converts string to UUID, returning zero UUID on error.
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// AlwaysAllowPermissionChecker is a development permission checker that always returns true.
// TODO: Replace with proper RBAC implementation in production.
type AlwaysAllowPermissionChecker struct{}

// NewAlwaysAllowPermissionChecker creates a new always-allow permission checker.
func NewAlwaysAllowPermissionChecker() *AlwaysAllowPermissionChecker {
	return &AlwaysAllowPermissionChecker{}
}

// Can always returns true for development purposes.
func (c *AlwaysAllowPermissionChecker) Can(ctx context.Context, actorID, tenantID string, action string) (bool, error) {
	return true, nil
}