package identity

import (
	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/app"
)

// RegisterRoutes registers identity routes.
func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {
	h := NewAuthHandler(NewIdentityService(
		NewIdentityRepositoryPostgres(deps.Postgres),
		deps.Config,
		deps.Logger,
		NewRBACPermissionChecker(deps.Postgres),
		deps.Postgres,
	), deps.Logger)
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		// TODO: Add more routes (logout, password reset, etc.) as needed
	}
}