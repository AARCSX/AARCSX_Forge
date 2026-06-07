package middleware

import (

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AARCSX/AARCSX_Forge/internal/app"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
)

// ContextMiddleware injects request-scoped values into the context:
// RequestID, TraceID, and optionally extracts TenantID from headers.
// This should be placed early in the middleware chain.
func ContextMiddleware(deps *app.RuntimeDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or extract request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Generate or extract trace ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Extract tenant ID from header if provided
		var tenantID string
		if tenantHeader := c.GetHeader("X-Tenant-ID"); tenantHeader != "" {
			tenantID = tenantHeader
		}

		// Create base context with request and trace IDs
		ctx := contextx.WithRequestID(c.Request.Context(), requestID)
		ctx = contextx.WithTraceID(ctx, traceID)

		// Add tenant ID if provided
		if tenantID != "" {
			ctx = contextx.WithTenantID(ctx, tenantID)
		}

		// Update request with new context
		c.Request = c.Request.WithContext(ctx)

		// Add context fields to logger for this request
		if deps.Logger != nil {
			c.Set("logger", deps.Logger.Sugar().With(
				"request_id", requestID,
				"trace_id", traceID,
				"tenant_id", tenantID,
			))
		}

		// Continue to next handler
		c.Next()
	}
}

// ContextValues returns a helper function to get logger with context values
// from gin.Context. Returns nil logger if context is not available.
func ContextValues(c *gin.Context) *zap.SugaredLogger {
	if logger, exists := c.Get("logger"); exists {
		if sugaredLogger, ok := logger.(*zap.SugaredLogger); ok {
			return sugaredLogger
		}
	}
	return nil
}

// ContextExtractor extracts common context values for use in handlers/services
type ContextExtractor struct{}

// NewContextExtractor creates a new context extractor
func NewContextExtractor() *ContextExtractor {
	return &ContextExtractor{}
}

// GetRequestID returns the request ID from context
func (e *ContextExtractor) GetRequestID(c *gin.Context) string {
	if ctx := c.Request.Context(); ctx != nil {
		return contextx.RequestID(ctx)
	}
	return ""
}

// GetTraceID returns the trace ID from context
func (e *ContextExtractor) GetTraceID(c *gin.Context) string {
	if ctx := c.Request.Context(); ctx != nil {
		return contextx.TraceID(ctx)
	}
	return ""
}

// GetTenantID returns the tenant ID from context
func (e *ContextExtractor) GetTenantID(c *gin.Context) string {
	if ctx := c.Request.Context(); ctx != nil {
		return contextx.TenantID(ctx)
	}
	return ""
}

// GetIdentity returns the identity context from context
func (e *ContextExtractor) GetIdentity(c *gin.Context) contextx.IdentityContext {
	if ctx := c.Request.Context(); ctx != nil {
		return contextx.Identity(ctx)
	}
	return contextx.IdentityContext{}
}