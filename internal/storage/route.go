package storage

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/app"
)

// RegisterRoutes registers storage routes.
func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {
	service, err := NewStorageService(&deps.Config, deps.Logger)
	if err != nil {
		// In a real application, you'd handle this error appropriately
		// For now, we'll panic during startup if storage service fails to initialize
		panic(fmt.Errorf("failed to initialize storage service: %w", err))
	}
	h := NewStorageHandler(service.WithRepository(NewObjectRepositoryPostgres(deps.Postgres)), deps.Logger)

	storage := router.Group("/storage")
	{
		storage.POST("/upload-url", h.GetUploadURL)
		storage.POST("/download-url", h.GetDownloadURL)
		// TODO: Add more routes (direct upload, metadata, etc.) as needed
	}
}