package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AARCSX/AARCSX_Forge/internal/logger"
)

// AuditLogService handles business logic for audit logs.
// Service layer only does business logic - no DB access or request handling.
type AuditLogService struct {
	repo        Repository
	logger      logger.Logger
	auditLogger Logger
	metrics     Metrics
	tracer      *Tracer
}

// Repository defines the interface for audit log persistence.
type Repository interface {
	Create(context.Context, AuditLog) (int64, error)
	GetByID(context.Context, int64) (*AuditLog, error)
	List(context.Context, int64, int64, string, string, int, int) ([]AuditLog, int64, error)
}

// NewAuditLogService creates a new audit log service.
func NewAuditLogService(repo Repository, logger logger.Logger, auditLogger Logger, metrics Metrics, tracer *Tracer) *AuditLogService {
	return &AuditLogService{
		repo:     repo,
		logger:   logger,
		auditLogger: auditLogger,
		metrics:  metrics,
		tracer:   tracer,
	}
}

// CreateAuditLog creates a new audit log entry after applying business rules.
func (s *AuditLogService) CreateAuditLog(ctx context.Context, input AuditLogCreateDTO) (AuditLogResponse, error) {
	// Start tracing span
	span, ctx := s.tracer.StartSpan(ctx, "CreateAuditLog", "")
	defer s.tracer.EndSpan(span)

	// Add span attributes
	s.tracer.SetSpanAttribute(span, "tenant_id", fmt.Sprint(input.TenantID))
	s.tracer.SetSpanAttribute(span, "user_id", fmt.Sprint(input.UserID))
	s.tracer.SetSpanAttribute(span, "action", input.Action)
	s.tracer.SetSpanAttribute(span, "resource", input.Resource)

	// Business logic: Validate input (beyond basic validation)
	if input.TenantID <= 0 {
		s.tracer.SetSpanStatus(span, "ERROR", "Invalid tenant ID")
		s.metrics.IncCounter("audit_log_create_validation_errors", map[string]string{"reason": "invalid_tenant_id"})
		return AuditLogResponse{}, ErrInvalidTenantID
	}
	if input.UserID <= 0 {
		s.tracer.SetSpanStatus(span, "ERROR", "Invalid user ID")
		s.metrics.IncCounter("audit_log_create_validation_errors", map[string]string{"reason": "invalid_user_id"})
		return AuditLogResponse{}, ErrInvalidUserID
	}
	if input.Action == "" {
		s.tracer.SetSpanStatus(span, "ERROR", "Empty action")
		s.metrics.IncCounter("audit_log_create_validation_errors", map[string]string{"reason": "empty_action"})
		return AuditLogResponse{}, ErrEmptyAction
	}
	if input.Resource == "" {
		s.tracer.SetSpanStatus(span, "ERROR", "Empty resource")
		s.metrics.IncCounter("audit_log_create_validation_errors", map[string]string{"reason": "empty_resource"})
		return AuditLogResponse{}, ErrEmptyResource
	}

	// Business logic: Prepare audit log entity
	auditLog := AuditLog{
		TenantID:   input.TenantID,
		UserID:     input.UserID,
		Action:     input.Action,
		Resource:   input.Resource,
		ResourceID: input.ResourceID,
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
		CreatedAt:  time.Now().UnixMilli(),
	}

	// Business logic: Handle details (serialize if needed)
	if input.Details != "" {
		// Validate that details is valid JSON if provided
		var js json.RawMessage
		if err := json.Unmarshal([]byte(input.Details), &js); err != nil {
			s.logger.Error(ctx, "Invalid JSON in audit log details", err, map[string]any{"details": input.Details})
			// Business rule: Store as-is if invalid JSON, but log warning
			s.tracer.AddSpanEvent(span, "details_validation_warning", map[string]string{
				"error":   err.Error(),
				"details": input.Details,
			})
		}
		auditLog.Details = input.Details
	}

	// Delegate to repository for persistence (no DB access in service)
	startTime := time.Now()
	id, err := s.repo.Create(ctx, auditLog)
	duration := time.Since(startTime).Seconds()

	if err != nil {
		s.tracer.SetSpanStatus(span, "ERROR", err.Error())
		s.metrics.IncCounter("audit_log_create_errors", map[string]string{"operation": "repository_create"})
		s.metrics.ObserveHistogram("audit_log_create_duration_seconds", duration, map[string]string{"status": "error"})
		return AuditLogResponse{}, err
	}

	auditLog.ID = id

	// Business logic: Convert to response DTO
	response := AuditLogResponse{
		ID:        auditLog.ID,
		TenantID:  auditLog.TenantID,
		UserID:    auditLog.UserID,
		Action:    auditLog.Action,
		Resource:  auditLog.Resource,
		ResourceID: auditLog.ResourceID,
		Details:   auditLog.Details,
		IPAddress: auditLog.IPAddress,
		UserAgent: auditLog.UserAgent,
		CreatedAt: auditLog.CreatedAt,
	}

	// Business logic: Fire audit event for real-time processing (if needed)
	s.auditLogger.Info(ctx, "Audit log created", map[string]any{
		"id":        id,
		"tenant_id": auditLog.TenantID,
		"user_id":   auditLog.UserID,
		"action":    auditLog.Action,
		"resource":  auditLog.Resource,
	})

	// Record success metrics
	s.tracer.SetSpanStatus(span, "OK", "Audit log created successfully")
	s.metrics.IncCounter("audit_logs_created_total", map[string]string{})
	s.metrics.ObserveHistogram("audit_log_create_duration_seconds", duration, map[string]string{"status": "success"})
	s.metrics.ObserveHistogram("audit_log_create_duration_seconds", duration, map[string]string{"action": input.Action})
	s.metrics.ObserveHistogram("audit_log_create_duration_seconds", duration, map[string]string{"resource": input.Resource})

	return response, nil
}

// GetAuditLog retrieves an audit log entry by ID after applying business rules.
func (s *AuditLogService) GetAuditLog(ctx context.Context, id int64) (*AuditLogResponse, error) {
	// Start tracing span
	span, ctx := s.tracer.StartSpan(ctx, "GetAuditLog", "")
	defer s.tracer.EndSpan(span)

	// Add span attributes
	s.tracer.SetSpanAttribute(span, "id", fmt.Sprint(id))

	// Business logic: Validate ID
	if id <= 0 {
		s.tracer.SetSpanStatus(span, "ERROR", "Invalid ID")
		s.metrics.IncCounter("audit_log_get_validation_errors", map[string]string{"reason": "invalid_id"})
		return nil, ErrInvalidID
	}

	// Delegate to repository for data access
	startTime := time.Now()
	log, err := s.repo.GetByID(ctx, id)
	duration := time.Since(startTime).Seconds()

	if err != nil {
		s.tracer.SetSpanStatus(span, "ERROR", err.Error())
		s.metrics.IncCounter("audit_log_get_errors", map[string]string{"operation": "repository_get_by_id"})
		s.metrics.ObserveHistogram("audit_log_get_duration_seconds", duration, map[string]string{"status": "error"})
		return nil, err
	}
	if log == nil {
		s.tracer.SetSpanStatus(span, "ERROR", "Audit log not found")
		s.metrics.IncCounter("audit_log_get_not_found", map[string]string{})
		s.metrics.ObserveHistogram("audit_log_get_duration_seconds", duration, map[string]string{"status": "not_found"})
		return nil, ErrNotFound
	}

	// Business logic: Convert to response DTO
	response := &AuditLogResponse{
		ID:        log.ID,
		TenantID:  log.TenantID,
		UserID:    log.UserID,
		Action:    log.Action,
		Resource:  log.Resource,
		ResourceID: log.ResourceID,
		Details:   log.Details,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
		CreatedAt: log.CreatedAt,
	}

	// Record success metrics
	s.tracer.SetSpanStatus(span, "OK", "Audit log retrieved successfully")
	s.metrics.IncCounter("audit_logs_retrieved_total", map[string]string{})
	s.metrics.ObserveHistogram("audit_log_get_duration_seconds", duration, map[string]string{"status": "success"})

	return response, nil
}

// ListAuditLogs retrieves paginated audit log entries with filtering.
func (s *AuditLogService) ListAuditLogs(ctx context.Context, tenantID int64, userID int64,
	action string, resource string, page int, pageSize int) (*AuditLogListResponse, error) {

	// Start tracing span
	span, ctx := s.tracer.StartSpan(ctx, "ListAuditLogs", "")
	defer s.tracer.EndSpan(span)

	// Add span attributes
	s.tracer.SetSpanAttribute(span, "tenant_id", fmt.Sprint(tenantID))
	s.tracer.SetSpanAttribute(span, "user_id", fmt.Sprint(userID))
	s.tracer.SetSpanAttribute(span, "action", action)
	s.tracer.SetSpanAttribute(span, "resource", resource)
	s.tracer.SetSpanAttribute(span, "page", fmt.Sprint(page))
	s.tracer.SetSpanAttribute(span, "page_size", fmt.Sprint(pageSize))

	// Business logic: Validate pagination
	if page <= 0 {
		page = 1
		s.tracer.AddSpanEvent(span, "page_validation_adjusted", map[string]string{"adjusted_to": "1"})
	}
	if pageSize <= 0 {
		pageSize = 20
		s.tracer.AddSpanEvent(span, "page_size_validation_adjusted", map[string]string{"adjusted_to": "20"})
	}
	if pageSize > 100 { // Business rule: Max page size
		pageSize = 100
		s.tracer.AddSpanEvent(span, "page_size_validation_adjusted", map[string]string{"adjusted_to": "100"})
	}

	// Delegate to repository for data access
	startTime := time.Now()
	logs, totalCount, err := s.repo.List(ctx, tenantID, userID, action, resource, page, pageSize)
	duration := time.Since(startTime).Seconds()

	if err != nil {
		s.tracer.SetSpanStatus(span, "ERROR", err.Error())
		s.metrics.IncCounter("audit_list_errors", map[string]string{"operation": "repository_list"})
		s.metrics.ObserveHistogram("audit_list_duration_seconds", duration, map[string]string{"status": "error"})
		return nil, err
	}

	// Business logic: Convert to response DTOs
	var items []AuditLogResponse
	for _, log := range logs {
		items = append(items, AuditLogResponse{
			ID:        log.ID,
			TenantID:  log.TenantID,
			UserID:    log.UserID,
			Action:    log.Action,
			Resource:  log.Resource,
			ResourceID: log.ResourceID,
			Details:   log.Details,
			IPAddress: log.IPAddress,
			UserAgent: log.UserAgent,
			CreatedAt: log.CreatedAt,
		})
	}

	// Business logic: Calculate pagination metadata
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	// Record success metrics
	s.tracer.SetSpanStatus(span, "OK", "Audit logs listed successfully")
	s.metrics.IncCounter("audit_logs_listed_total", map[string]string{})
	s.metrics.ObserveHistogram("audit_list_duration_seconds", duration, map[string]string{"status": "success"})
	s.metrics.ObserveHistogram("audit_list_items_count", float64(len(items)), map[string]string{})
	s.metrics.SetGauge("audit_log_total_count", float64(totalCount), map[string]string{})

	return &AuditLogListResponse{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Error types for business logic validation
var (
	ErrInvalidTenantID = NewValidationError("tenant ID must be positive")
	ErrInvalidUserID   = NewValidationError("user ID must be positive")
	ErrEmptyAction     = NewValidationError("action cannot be empty")
	ErrEmptyResource   = NewValidationError("resource cannot be empty")
	ErrInvalidID       = NewValidationError("ID must be positive")
	ErrNotFound        = NewNotFoundError("audit log not found")
)

// ValidationError represents a validation error in business logic.
type ValidationError struct {
	Message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func (e *ValidationError) Error() string {
	return "validation error: " + e.Message
}

// NotFoundError represents a not found error in business logic.
type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

func (e *NotFoundError) Error() string {
	return "not found: " + e.Message
}