package observability

// AuditLog represents an audit log entry stored in the database.
type AuditLog struct {
	ID        int64  `json:"id" db:"id"`
	TenantID  int64  `json:"tenant_id" db:"tenant_id"`
	UserID    int64  `json:"user_id" db:"user_id"`
	Action    string `json:"action" db:"action"`
	Resource  string `json:"resource" db:"resource"`
	ResourceID string `json:"resource_id" db:"resource_id"`
	Details   string `json:"details" db:"details"` // JSON stringified details
	IPAddress string `json:"ip_address" db:"ip_address"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}