package observability

import (
	"github.com/gin-gonic/gin"

	"github.com/AARCSX/AARCSX_Forge/internal/database"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
)

// RegisterRoutes registers observability routes.
func RegisterRoutes(router *gin.RouterGroup, postgres *database.Postgres, logger *logger.Logger, metrics Metrics, tracer *Tracer) {
	h := NewAuditLogHandler(
		NewAuditLogService(
			NewAuditLogRepository(postgres, logger),
			*logger,
			logger, // Using same logger for audit logging for now
			metrics,
			tracer,
		),
		logger,
		metrics,
		tracer,
	)
	auditLogs := router.Group("/audit-logs")
	{
		auditLogs.POST("", h.CreateAuditLog)
		auditLogs.GET("/:id", h.GetAuditLog)
		auditLogs.GET("", h.ListAuditLogs)
		// TODO: Add more routes (export, delete, etc.) as needed
	}
}