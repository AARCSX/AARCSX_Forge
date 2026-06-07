package observability

// AuditLogCreateDTO represents the data needed to create an audit log entry.
type AuditLogCreateDTO struct {
	TenantID   int64  `json:"tenant_id" validate:"required"`
	UserID     int64  `json:"user_id" validate:"required"`
	Action     string `json:"action" validate:"required"`
	Resource   string `json:"resource" validate:"required"`
	ResourceID string `json:"resource_id,omitempty"`
	Details    string `json:"details,omitempty"` // JSON stringified details
	IPAddress  string `json:"ip_address,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
}

// AuditLogResponse represents the audit log entry for API responses.
type AuditLogResponse struct {
	ID        int64  `json:"id"`
	TenantID  int64  `json:"tenant_id"`
	UserID    int64  `json:"user_id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	ResourceID string `json:"resource_id"`
	Details   string `json:"details"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	CreatedAt int64  `json:"created_at"`
}

// AuditLogListResponse represents a paginated list of audit logs.
type AuditLogListResponse struct {
	Items      []AuditLogResponse `json:"items"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}