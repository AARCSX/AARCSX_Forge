package observability

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
)

// AuditLogHandler handles HTTP requests for audit logs.
// Handler layer only processes requests - no business logic or DB access.
type AuditLogHandler struct {
	service Service
	logger  *logger.Logger
	metrics Metrics
	tracer  *Tracer
}

// Service defines the interface for audit log business logic.
type Service interface {
	CreateAuditLog(context.Context, AuditLogCreateDTO) (AuditLogResponse, error)
	GetAuditLog(context.Context, int64) (*AuditLogResponse, error)
	ListAuditLogs(context.Context, int64, int64, string, string, int, int) (*AuditLogListResponse, error)
}

// NewAuditLogHandler creates a new audit log handler.
func NewAuditLogHandler(service Service, logger *logger.Logger, metrics Metrics, tracer *Tracer) *AuditLogHandler {
	return &AuditLogHandler{service: service, logger: logger, metrics: metrics, tracer: tracer}
}

// CreateAuditLog handles POST /audit-logs
func (h *AuditLogHandler) CreateAuditLog(c *gin.Context) {
	// Start tracing span
	span, ctx := h.tracer.StartSpan(c.Request.Context(), "HTTP_CREATE_AUDIT_LOG", "")
	defer h.tracer.EndSpan(span)

	var input AuditLogCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		h.tracer.SetSpanStatus(span, "ERROR", "Invalid request")
		h.metrics.IncCounter("http_requests_total", map[string]string{
			"method": "POST",
			"endpoint": "/audit-logs",
			"status": "400",
			"reason": "invalid_request",
		})
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.CreateAuditLog(ctx, input)
	if err != nil {
		h.handleServiceError(c, err)
		h.tracer.SetSpanStatus(span, "ERROR", err.Error())
		h.metrics.IncCounter("http_requests_total", map[string]string{
			"method": "POST",
			"endpoint": "/audit-logs",
			"status": "500",
			"reason": "internal_error",
		})
		return
	}

	// Record success metrics
	h.tracer.SetSpanStatus(span, "OK", "Audit log created")
	h.metrics.IncCounter("http_requests_total", map[string]string{
		"method": "POST",
		"endpoint": "/audit-logs",
		"status": "201",
	})
	h.metrics.ObserveHistogram("http_request_duration_seconds", 0.1, map[string]string{ // Placeholder duration
		"method": "POST",
		"endpoint": "/audit-logs",
		"status": "201",
	})

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// GetAuditLog handles GET /audit-logs/:id
func (h *AuditLogHandler) GetAuditLog(c *gin.Context) {
	// Start tracing span
	span, ctx := h.tracer.StartSpan(c.Request.Context(), "HTTP_GET_AUDIT_LOG", "")
	defer h.tracer.EndSpan(span)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Sugar().Errorw("invalid ID parameter", "error", err, "param", c.Param("id"))
		h.tracer.SetSpanStatus(span, "ERROR", "Invalid ID parameter")
		h.metrics.IncCounter("http_requests_total", map[string]string{
			"method": "GET",
			"endpoint": "/audit-logs/:id",
			"status": "400",
			"reason": "invalid_id",
		})
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_ID",
			Message: "invalid ID",
		})
		return
	}

	// Add span attribute for ID
	h.tracer.SetSpanAttribute(span, "audit_log_id", c.Param("id"))

	res, err := h.service.GetAuditLog(ctx, id)
	if err != nil {
		h.handleServiceError(c, err)
		h.tracer.SetSpanStatus(span, "ERROR", err.Error())
		statusCode := "500"
		reason := "internal_error"
		switch err.(type) {
		case *ValidationError:
			statusCode = "400"
			reason = "validation_error"
		case *NotFoundError:
			statusCode = "404"
			reason = "not_found"
		}
		h.metrics.IncCounter("http_requests_total", map[string]string{
			"method": "GET",
			"endpoint": "/audit-logs/:id",
			"status": statusCode,
			"reason": reason,
		})
		return
	}

	// Record success metrics
	h.tracer.SetSpanStatus(span, "OK", "Audit log retrieved")
	h.metrics.IncCounter("http_requests_total", map[string]string{
		"method": "GET",
		"endpoint": "/audit-logs/:id",
		"status": "200",
	})
	h.metrics.ObserveHistogram("http_request_duration_seconds", 0.1, map[string]string{ // Placeholder duration
		"method": "GET",
		"endpoint": "/audit-logs/:id",
		"status": "200",
	})

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// ListAuditLogs handles GET /audit-logs
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	// Start tracing span
	span, ctx := h.tracer.StartSpan(c.Request.Context(), "HTTP_LIST_AUDIT_LOGS", "")
	defer h.tracer.EndSpan(span)

	// Parse query parameters
	tenantID, _ := strconv.ParseInt(c.Query("tenant_id"), 10, 64)
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	action := c.Query("action")
	resource := c.Query("resource")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Add span attributes for filtering
	if tenantID != 0 {
		h.tracer.SetSpanAttribute(span, "filter_tenant_id", c.Query("tenant_id"))
	}
	if userID != 0 {
		h.tracer.SetSpanAttribute(span, "filter_user_id", c.Query("user_id"))
	}
	if action != "" {
		h.tracer.SetSpanAttribute(span, "filter_action", action)
	}
	if resource != "" {
		h.tracer.SetSpanAttribute(span, "filter_resource", resource)
	}
	h.tracer.SetSpanAttribute(span, "page", c.DefaultQuery("page", "1"))
	h.tracer.SetSpanAttribute(span, "page_size", c.DefaultQuery("page_size", "20"))

	res, err := h.service.ListAuditLogs(ctx, tenantID, userID, action, resource, page, pageSize)
	if err != nil {
		h.handleServiceError(c, err)
		h.tracer.SetSpanStatus(span, "ERROR", err.Error())
		statusCode := "500"
		reason := "internal_error"
		switch err.(type) {
		case *ValidationError:
			statusCode = "400"
			reason = "validation_error"
		case *NotFoundError:
			statusCode = "404"
			reason = "not_found"
		}
		h.metrics.IncCounter("http_requests_total", map[string]string{
			"method": "GET",
			"endpoint": "/audit-logs",
			"status": statusCode,
			"reason": reason,
		})
		return
	}

	// Record success metrics
	h.tracer.SetSpanStatus(span, "OK", "Audit logs listed")
	h.metrics.IncCounter("http_requests_total", map[string]string{
		"method": "GET",
		"endpoint": "/audit-logs",
		"status": "200",
	})
	h.metrics.ObserveHistogram("http_request_duration_seconds", 0.1, map[string]string{ // Placeholder duration
		"method": "GET",
		"endpoint": "/audit-logs",
		"status": "200",
	})
	h.metrics.ObserveHistogram("http_request_items_count", float64(len(res.Items)), map[string]string{
		"method": "GET",
		"endpoint": "/audit-logs",
	})

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// handleServiceError converts service errors to appropriate HTTP responses.
func (h *AuditLogHandler) handleServiceError(c *gin.Context, err error) {
	switch err.(type) {
	case *ValidationError:
		h.logger.Sugar().Warnw("validation error", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
	case *NotFoundError:
		h.logger.Sugar().Warnw("not found", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusNotFound, httpx.ErrorDetail{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
	default:
		h.logger.Sugar().Errorw("internal error", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		})
	}
}