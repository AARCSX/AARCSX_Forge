package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/app"
)

// Handler handles health check requests.
type Handler struct {
	Deps *app.RuntimeDeps
}

// NewHandler creates a new health check handler.
func NewHandler(deps *app.RuntimeDeps) *Handler {
	return &Handler{Deps: deps}
}

// RegisterRoutes registers health check routes.
func RegisterRoutes(router *gin.RouterGroup, deps *app.RuntimeDeps) {
	h := NewHandler(deps)
	health := router.Group("/health")
	{
		health.GET("", h.Check)
		health.GET("/ready", h.Ready)
		health.GET("/live", h.Live)
	}
}

// Check returns overall service health.
func (h *Handler) Check(c *gin.Context) {
	// Check database connection
	dbErr := h.Deps.Postgres.DB.PingContext(c.Request.Context())
	// Check redis connection
	redisErr := h.Deps.Redis.Client.Ping(c.Request.Context()).Err()

	if dbErr != nil || redisErr != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"database": dbErr == nil,
			"redis":    redisErr == nil,
			"errors": gin.H{
				"database": mapError(dbErr),
				"redis":    mapError(redisErr),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"database": true,
		"redis":    true,
	})
}

// Ready checks if the service is ready to serve traffic.
func (h *Handler) Ready(c *gin.Context) {
	h.Check(c)
}

// Live checks if the service is alive.
func (h *Handler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func mapError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}