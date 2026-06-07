package tenants

import (
	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/app"
)

// RegisterRoutes registers tenant routes.
func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {
	h := NewTenantHandler(NewTenantService(NewTenantRepositoryPostgres(deps.Postgres), deps.Logger), deps.Logger)
	tenants := router.Group("/tenants")
	{
		tenants.POST("", h.CreateTenant)
		tenants.GET("/:id", h.GetTenant)
		tenants.PUT("/:id", h.UpdateTenant)
		// TODO: Add more routes (list, delete, etc.) as needed
	}
}