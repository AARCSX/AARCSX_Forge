package httpx

import (
	"net/http"
)

type ResponseEnvelope struct {
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	TraceID string       `json:"trace_id"`
}

type ErrorDetail struct {
	Code           string         `json:"code"`
	Message        string         `json:"message"`
	LocalizationID string         `json:"localization_id,omitempty"`
	Meta           map[string]any `json:"meta,omitempty"`
}

type JSONWriter interface {
	JSON(statusCode int, payload any)
}

func Success(w JSONWriter, traceID string, data any) {
	w.JSON(http.StatusOK, ResponseEnvelope{
		Success: true,
		Data:    data,
		TraceID: traceID,
	})
}

func Failure(w JSONWriter, statusCode int, traceID string, err ErrorDetail) {
	w.JSON(statusCode, ResponseEnvelope{
		Success: false,
		Error:   &err,
		TraceID: traceID,
	})
}
