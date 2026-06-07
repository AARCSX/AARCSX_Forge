package storage

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
)

// StorageHandler handles HTTP requests for storage operations.
type StorageHandler struct {
	service Service
	logger  *logger.Logger
}

// NewStorageHandler creates a new storage handler.
func NewStorageHandler(service Service, logger *logger.Logger) *StorageHandler {
	return &StorageHandler{service: service, logger: logger}
}

// GetUploadURL handles POST /storage/upload-url
func (h *StorageHandler) GetUploadURL(c *gin.Context) {
	var req CreateUploadURLInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.CreateUploadURL(c.Request.Context(), req)
	if err != nil {
		h.logger.Sugar().Errorw("failed to create upload URL", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create upload URL",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// GetDownloadURL handles POST /storage/download-url
func (h *StorageHandler) GetDownloadURL(c *gin.Context) {
	var req CreateDownloadURLInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.CreateDownloadURL(c.Request.Context(), req)
	if err != nil {
		h.logger.Sugar().Errorw("failed to create download URL", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create download URL",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}