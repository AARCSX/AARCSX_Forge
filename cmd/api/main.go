package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AARCSX/AARCSX_Forge/internal/app"
	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/AARCSX/AARCSX_Forge/internal/health"
	"github.com/AARCSX/AARCSX_Forge/internal/identity"
	"github.com/AARCSX/AARCSX_Forge/internal/notifications"
	"github.com/AARCSX/AARCSX_Forge/internal/observability"
	contextMiddleware "github.com/AARCSX/AARCSX_Forge/internal/platform/contextx/middleware"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
	platformMiddleware "github.com/AARCSX/AARCSX_Forge/internal/platform/middleware"
	"github.com/AARCSX/AARCSX_Forge/internal/storage"
	"github.com/AARCSX/AARCSX_Forge/internal/tenants"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Build runtime dependencies
	deps := app.BuildRuntimeDeps(cfg)
	defer func() {
		// Ensure shutdown is called
		if err := deps.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	// Set up Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Register health routes
	health.RegisterRoutes(r.Group("/"), &deps)

	// Setup context middleware FIRST - sets requestID, traceID, tenantID
	contextMiddleware := contextMiddleware.ContextMiddleware(&deps)

	// Define API routes (to be expanded in later sprints)
	api := r.Group("/api/v1")
	{
		// Placeholder for future routes
		api.GET("/ping", func(c *gin.Context) {
			httpx.Success(httpx.NewGinResponseWriter(c), gin.H{
				"message": "pong",
				"env":     deps.Config.Env,
			})
		})
	}

	// Setup tenant middleware
	tenantMiddleware := platformMiddleware.TenantMiddleware(tenants.NewTenantService(tenants.NewTenantRepositoryPostgres(deps.Postgres), deps.Logger))

	// Auth middleware - validates JWT and sets identity context
	authMiddleware := platformMiddleware.AuthMiddleware(&deps)

	// Tenant-scoped routes (need context, tenant, and auth)
	tenantsGroup := api.Group("/tenants")
	tenantsGroup.Use(contextMiddleware)
	tenantsGroup.Use(tenantMiddleware)
	tenantsGroup.Use(authMiddleware)
	{
		// Register tenant routes
		tenants.RegisterRoutes(tenantsGroup, &deps)
	}

	// Auth routes (public, no tenant required, but need context for logging/tracing)
	authGroup := api.Group("/auth")
	authGroup.Use(contextMiddleware)
	{
		// Register auth routes
		identity.RegisterRoutes(authGroup, &deps)
	}

	// Notification routes (tenant-scoped)
	notificationGroup := api.Group("/notifications")
	notificationGroup.Use(contextMiddleware)
	notificationGroup.Use(tenantMiddleware)
	notificationGroup.Use(authMiddleware)
	{
		// Register notification routes
		notifications.RegisterRoutes(notificationGroup, &deps)
	}

	// Storage routes (tenant-scoped)
	storageGroup := api.Group("/storage")
	storageGroup.Use(contextMiddleware)
	storageGroup.Use(tenantMiddleware)
	storageGroup.Use(authMiddleware)
	{
		// Register storage routes
		storage.RegisterRoutes(storageGroup, &deps)
	}

	// Observability routes (tenant-scoped)
	observabilityGroup := api.Group("/observability")
	observabilityGroup.Use(contextMiddleware)
	observabilityGroup.Use(tenantMiddleware)
	observabilityGroup.Use(authMiddleware)
	{
		// Register observability routes
		observability.RegisterRoutes(observabilityGroup, deps.Postgres, deps.Logger, deps.Metrics, deps.Tracer)
	}

	// Start server
	srv := &http.Server{
		Addr:    deps.Config.API.HTTPAddress,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	// Call shutdown for dependencies (already deferred, but we do it again for clarity)
	if err := deps.Shutdown(ctx); err != nil {
		log.Printf("dependency shutdown error: %v", err)
	}

	log.Println("server exiting")
}