package notifications

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/AARCSX/AARCSX_Forge/internal/app"
)

// RegisterRoutes registers notification routes.
func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {
	service, err := NewNotificationService(deps.Logger)
	if err != nil {
		// In a real application, you'd handle this error appropriately
		// For now, we'll panic during startup if notification service fails to initialize
		panic(fmt.Errorf("failed to initialize notification service: %w", err))
	}
	h := NewNotificationHandler(service.WithRepository(NewDeliveryJobPostgres(deps.Postgres)), deps.Logger)

	notifications := router.Group("/notifications")
	{
		notifications.POST("/email", h.SendEmail)
		// TODO: Add more routes (template management, delivery status, etc.) as needed
	}
}