package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware returns a gin.HandlerFunc that limits the request rate.
func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	// Create a rate limiter per client IP.
	limiters := make(map[string]*rate.Limiter)
	var mutex sync.Mutex

	return func(c *gin.Context) {
		// Get the client IP.
		ip := c.ClientIP()
		mutex.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			limiters[ip] = limiter
		}
		mutex.Unlock()

		// Check if the request is allowed.
		if !limiter.Allow() {
			c.AbortWithStatusJSON(httpTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		// Continue to next handler.
		c.Next()
	}
}

// Constants for HTTP status codes.
const (
	httpTooManyRequests = http.StatusTooManyRequests // 429
)