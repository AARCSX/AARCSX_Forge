package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/AARCSX/AARCSX_Forge/internal/app"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
	"github.com/AARCSX/AARCSX_Forge/internal/identity"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
)

// AuthMiddleware validates JWT tokens and sets user identity context.
// Assumes ContextMiddleware has already set requestID, traceID, and tenantID.
func AuthMiddleware(deps *app.RuntimeDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		// Check for Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(deps.Config.JWT.SigningKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// Extract user ID from subject
		userIDStr, ok := claims["sub"].(string)
		if !ok || userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			return
		}

		// Get user from database
		user, err := identity.NewIdentityRepositoryPostgres(deps.Postgres).GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		// Check if user is active
		if user.Status != "active" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user account is not active"})
			return
		}

		// Get user roles and permissions
		roles, permissions, err := getUserRolesAndPermissions(c.Request.Context(), deps.Postgres, userID, user.TenantID)
		if err != nil {
			deps.Logger.Sugar().Errorw("failed to get user roles and permissions", "error", err, "user_id", userID)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user permissions"})
			return
		}

		// Create identity context
		identityCtx := contextx.IdentityContext{
			ActorID:     userID.String(),
			Roles:       roles,
			Permissions: permissions,
		}

		// Set identity context (requestID, traceID, tenantID should already be set by ContextMiddleware)
		ctx := contextx.WithIdentity(c.Request.Context(), identityCtx)

		// Update request with new context
		c.Request = c.Request.WithContext(ctx)

		// Continue to next handler
		c.Next()
	}
}

// getUserRolesAndPermissions retrieves roles and permissions for a user.
func getUserRolesAndPermissions(ctx context.Context, db *database.Postgres, userID, tenantID uuid.UUID) ([]string, []string, error) {
	var roles []string
	var permissions []string

	// Query to get roles and their permissions for the user in the specific tenant
	query := `
		SELECT r.name, r.permissions
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.tenant_id = $2
	`

	rows, err := db.DB.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Map to collect unique permissions
	permissionMap := make(map[string]bool)

	for rows.Next() {
		var roleName string
		var perms []byte // JSONB permissions

		if err := rows.Scan(&roleName, &perms); err != nil {
			return nil, nil, err
		}

		roles = append(roles, roleName)

		// Parse JSON permissions array
		var permList []string
		if err := db.DB.QueryRowContext(ctx, "SELECT jsonb_array_elements_text($1::jsonb)", perms).Scan(&permList); err != nil {
			// Try alternative approach for parsing JSONB array
			if err := db.DB.QueryRowContext(ctx, "SELECT COALESCE(array_to_json(ARRAY(SELECT jsonb_array_elements_text($1))), '[]'::jsonb)", perms).Scan(&permList); err != nil {
				// If we can't parse, continue with empty permissions
				continue
			}
		}

		for _, perm := range permList {
			permissionMap[perm] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	// Convert permission map to slice
	for perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return roles, permissions, nil
}

// PermissionMiddleware checks if the user has the required permission.
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get identity from context
		identityCtx := contextx.Identity(c.Request.Context())

		// Check if user has the required permission
		hasPermission := false
		for _, perm := range identityCtx.Permissions {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		// Continue to next handler
		c.Next()
	}
}

// RoleMiddleware checks if the user has one of the required roles.
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get identity from context
		identityCtx := contextx.Identity(c.Request.Context())

		// Check if user has any of the required roles
		hasRole := false
		userRoles := make(map[string]bool)
		for _, role := range identityCtx.Roles {
			userRoles[role] = true
		}

		for _, requiredRole := range requiredRoles {
			if userRoles[requiredRole] {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role permissions"})
			return
		}

		// Continue to next handler
		c.Next()
	}
}