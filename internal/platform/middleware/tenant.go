package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
	tenants "github.com/AARCSX/AARCSX_Forge/internal/tenants"
)

// TenantMiddleware resolves tenant from JWT or X-Tenant-ID header.
func TenantMiddleware(ts *tenants.TenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID string

		// 1. Check for tenant in JWT claims (if we had a middleware that populates context with JWT claims)
		// For now, we'll rely on the X-Tenant-ID header as the source of truth for tenant resolution.
		// In a real implementation, we would extract tenant from JWT and then verify against header or allow header to override.

		// For Sprint 2, we'll assume the tenant is passed via X-Tenant-ID header.
		if tenantHeader := c.GetHeader("X-Tenant-ID"); tenantHeader != "" {
			tenantID = tenantHeader
		} else {
			// If no header, we could try to get from JWT (but we don't have JWT middleware yet)
			// For now, return an error if no tenant is provided.
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Tenant-ID header"})
			return
		}

		// Validate tenant ID format (optional, but we can check if it's a valid UUID)
		if _, err := uuid.Parse(tenantID); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
			return
		}

		// Optionally, we can verify the tenant exists in the database (but this adds overhead).
		// For performance, we might skip this and rely on foreign key constraints in the DB.
		// We'll do a quick check in development, but in production we might skip.
		// For now, we'll skip the DB check to keep the middleware fast.

		// Store tenant ID in context
		ctx := contextx.WithTenantID(c.Request.Context(), tenantID)
		c.Request = c.Request.WithContext(ctx)

		// Continue to next handler
		c.Next()
	}
}