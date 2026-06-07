-- Add object metadata table for storage service
CREATE TABLE IF NOT EXISTS object_metadata (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    object_key VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for object metadata
CREATE INDEX IF NOT EXISTS idx_object_metadata_tenant_id ON object_metadata(tenant_id);
CREATE INDEX IF NOT EXISTS idx_object_metadata_object_key ON object_metadata(object_key);