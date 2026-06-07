package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
)

// ResponseEnvelope represents a standardized API response.
type ResponseEnvelope struct {
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	TraceID string       `json:"trace_id"`
}

// ErrorDetail represents standardized error information.
type ErrorDetail struct {
	Code           string         `json:"code"`
	Message        string         `json:"message"`
	LocalizationID string         `json:"localization_id,omitempty"`
	Meta           map[string]any `json:"meta,omitempty"`
}

// JSONWriter interface abstracts the JSON response writer.
type JSONWriter interface {
	JSON(statusCode int, payload any)
}

// GinResponseWriter implements JSONWriter for gin.Context.
type GinResponseWriter struct {
	c *gin.Context
}

// NewGinResponseWriter creates a new GinResponseWriter.
func NewGinResponseWriter(c *gin.Context) *GinResponseWriter {
	return &GinResponseWriter{c: c}
}

// JSON implements the JSONWriter interface for gin.Context.
func (w *GinResponseWriter) JSON(statusCode int, payload any) {
	w.c.JSON(statusCode, payload)
}

// Success returns a successful standardized response.
// Automatically extracts traceID from context if not provided.
func Success(w JSONWriter, data any) {
	// Try to get traceID from context if writer is GinResponseWriter
	traceID := ""
	if gw, ok := w.(*GinResponseWriter); ok && gw.c != nil {
		traceID = contextx.TraceID(gw.c.Request.Context())
	}

	w.JSON(http.StatusOK, ResponseEnvelope{
		Success: true,
		Data:    data,
		TraceID: traceID,
	})
}

// Failure returns an error standardized response.
// Automatically extracts traceID from context if not provided.
func Failure(w JSONWriter, statusCode int, err ErrorDetail) {
	// Try to get traceID from context if writer is GinResponseWriter
	traceID := ""
	if gw, ok := w.(*GinResponseWriter); ok && gw.c != nil {
		traceID = contextx.TraceID(gw.c.Request.Context())
	}

	w.JSON(statusCode, ResponseEnvelope{
		Success: false,
		Error:   &err,
		TraceID: traceID,
	})
}