package tenants

// CreateTenantRequest represents the request to create a tenant.
type CreateTenantRequest struct {
	Name  string `json:"name" validate:"required"`
	Slug  string `json:"slug" validate:"required,slug"`
	Status string `json:"status,omitempty"`
}

// UpdateTenantRequest represents the request to update a tenant.
type UpdateTenantRequest struct {
	Name  string `json:"name,omitempty"`
	Slug  string `json:"slug,omitempty,slug"`
	Status string `json:"status,omitempty"`
}

// TenantResponse represents the response for tenant operations.
type TenantResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}