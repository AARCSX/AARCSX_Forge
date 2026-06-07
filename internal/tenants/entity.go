package tenants

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system.
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	Settings  []byte    `json:"settings,omitempty"` // JSONB stored as []byte
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}