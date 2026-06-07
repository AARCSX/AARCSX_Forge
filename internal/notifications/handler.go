package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
)

// NotificationHandler handles HTTP requests for notification operations.
type NotificationHandler struct {
	service Service
	logger  *logger.Logger
}

// NewNotificationHandler creates a new notification handler.
func NewNotificationHandler(service Service, logger *logger.Logger) *NotificationHandler {
	return &NotificationHandler{service: service, logger: logger}
}

// SendEmail handles POST /notifications/email
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	var req QueueEmailInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	ctx := c.Request.Context()
	res, err := h.service.QueueEmail(ctx, req)
	if err != nil {
		h.logger.Sugar().Errorw("failed to queue email", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to queue email",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}