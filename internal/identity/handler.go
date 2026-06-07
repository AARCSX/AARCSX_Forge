package identity

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/httpx"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	service Service
	logger  *logger.Logger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(service Service, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{service: service, logger: logger}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		if err == ErrUserAlreadyExists {
			httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusConflict, httpx.ErrorDetail{
				Code:    "USER_ALREADY_EXISTS",
				Message: "user already exists",
			})
			return
		}
		h.logger.Sugar().Errorw("failed to register user", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to register user",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthenticateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.Authenticate(c.Request.Context(), &req)
	if err != nil {
		if err == ErrInvalidCredentials {
			httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusUnauthorized, httpx.ErrorDetail{
				Code:    "INVALID_CREDENTIALS",
				Message: "invalid credentials",
			})
			return
		}
		h.logger.Sugar().Errorw("failed to authenticate", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to authenticate",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}

// Refresh handles POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshTokenInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Sugar().Errorw("invalid request", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusBadRequest, httpx.ErrorDetail{
			Code:    "INVALID_REQUEST",
			Message: "invalid request",
		})
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if err == ErrInvalidCredentials {
			httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusUnauthorized, httpx.ErrorDetail{
				Code:    "INVALID_CREDENTIALS",
				Message: "invalid refresh token",
			})
			return
		}
		h.logger.Sugar().Errorw("failed to refresh token", "error", err)
		httpx.Failure(httpx.NewGinResponseWriter(c), http.StatusInternalServerError, httpx.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "failed to refresh token",
		})
		return
	}

	httpx.Success(httpx.NewGinResponseWriter(c), res)
}