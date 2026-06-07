package tenants

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
)

// TenantHandler handles HTTP requests for tenants.
type TenantHandler struct {
	service *TenantService
	logger  *logger.Logger
}

// NewTenantHandler creates a new tenant handler.
func NewTenantHandler(service *TenantService, logger *logger.Logger) *TenantHandler {
	return &TenantHandler{service: service, logger: logger}
}

// CreateTenant handles POST /tenants
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		h.logger.Sugar().Errorw("failed to create tenant", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create tenant",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// GetTenant handles GET /tenants/:id
func (h *TenantHandler) GetTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Sugar().Errorw("invalid tenant ID", "error", err, "id", idStr)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_TENANT_ID",
			Message: "invalid tenant ID",
		})
		return
	}

	res, err := h.service.GetTenantByID(c.Request.Context(), id)
	if err != nil {
		if err == ErrTenantNotFound {
			httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusNotFound, httpx.ErrorDetail{
				Code:    "TENANT_NOT_FOUND",
				Message: "tenant not found",
			})
			return
		}
		h.logger.Sugar().Errorw("failed to get tenant", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to get tenant",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// UpdateTenant handles PUT /tenants/:id
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Sugar().Errorw("invalid tenant ID", "error", err, "id", idStr)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_TENANT_ID",
			Message: "invalid tenant ID",
		})
		return
	}

	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.UpdateTenant(c.Request.Context(), id, &req)
	if err != nil {
		if err == ErrTenantNotFound {
			httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusNotFound, httpx.ErrorDetail{
				Code:    "TENANT_NOT_FOUND",
				Message: "tenant not found",
			})
			return
		}
		h.logger.Sugar().Errorw("failed to update tenant", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to update tenant",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}