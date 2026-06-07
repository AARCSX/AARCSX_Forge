-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    action VARCHAR(255) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    resource_id VARCHAR(255),
    details TEXT, -- JSON stringified details
    ip_address INET,
    user_agent TEXT,
    created_at BIGINT NOT NULL, -- Unix timestamp in milliseconds

    -- Indexes for common query patterns
    INDEX idx_audit_logs_tenant_id (tenant_id),
    INDEX idx_audit_logs_user_id (user_id),
    INDEX idx_audit_logs_action (action),
    INDEX idx_audit_logs_resource (resource),
    INDEX idx_audit_logs_created_at (created_at),
    INDEX idx_audit_logs_tenant_user (tenant_id, user_id),
    INDEX idx_audit_logs_tenant_action (tenant_id, action),
    INDEX idx_audit_logs_tenant_resource (tenant_id, resource)
);

-- Add comment on table
COMMENT ON TABLE audit_logs IS 'Audit log entries for tracking user actions and system events';

-- Add comments on columns
COMMENT ON COLUMN audit_logs.id IS 'Unique identifier for the audit log entry';
COMMENT ON COLUMN audit_logs.tenant_id IS 'ID of the tenant the action belongs to';
COMMENT ON COLUMN audit_logs.user_id IS 'ID of the user who performed the action';
COMMENT ON COLUMN audit_logs.action IS 'Action performed (e.g., CREATE_USER, UPDATE_TENANT)';
COMMENT ON COLUMN audit_logs.resource IS 'Resource that was acted upon (e.g., USER, TENANT, FILE)';
COMMENT ON COLUMN audit_logs.resource_id IS 'ID of the specific resource that was acted upon';
COMMENT ON COLUMN audit_logs.details IS 'Additional details about the action in JSON format';
COMMENT ON COLUMN audit_logs.ip_address IS 'IP address of the user who performed the action';
COMMENT ON COLUMN audit_logs.user_agent IS 'User agent string of the client';
COMMENT ON COLUMN audit_logs.created_at IS 'Timestamp when the action occurred (Unix milliseconds)';