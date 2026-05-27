package tenants

import "context"

type Service interface {
	Resolve(ctx context.Context, in ResolveInput) (TenantContext, error)
	Provision(ctx context.Context, in ProvisionInput) (Tenant, error)
}

type Repository interface {
	GetByID(ctx context.Context, tenantID string) (Tenant, error)
	GetMembership(ctx context.Context, tenantID, userID string) (Membership, error)
	Create(ctx context.Context, tenant Tenant) (Tenant, error)
}

type ResolveInput struct {
	JWTTenantID string
	HeaderValue string
}

type ProvisionInput struct {
	Name    string
	OwnerID string
}

type Tenant struct {
	ID   string
	Name string
}

type Membership struct {
	TenantID string
	UserID   string
	Roles    []string
}

type TenantContext struct {
	TenantID string
	Source   string
}
