-- Create notification_deliveries table
CREATE TABLE IF NOT EXISTS notification_deliveries (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for notification_deliveries
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_tenant_id ON notification_deliveries(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_status ON notification_deliveries(status);