package observability

import (
	"context"
	"time"

	"github.com/AARCSX/AARCSX_Forge/internal/platform/contextx"
	"github.com/google/uuid"
)

// generateID returns a new UUID string
func generateID() string {
	return uuid.New().String()
}

// TraceIDGenerator generates trace IDs for distributed tracing.
type TraceIDGenerator interface {
	GenerateTraceID() string
	GenerateSpanID() string
}

// SimpleTraceIDGenerator generates simple UUID-based trace IDs.
// In production, this would integrate with OpenTelemetry, Jaeger, etc.
type SimpleTraceIDGenerator struct{}

// NewSimpleTraceIDGenerator creates a new simple trace ID generator.
func NewSimpleTraceIDGenerator() *SimpleTraceIDGenerator {
	return &SimpleTraceIDGenerator{}
}

// GenerateTraceID generates a trace ID from context if available, otherwise creates new one.
func (g *SimpleTraceIDGenerator) GenerateTraceID() string {
	// Note: This method signature doesn't accept context, so we always generate new
	// In a real implementation with OpenTelemetry, this would extract from context
	return "trace-" + generateID()
}

// GenerateSpanID generates a span ID.
func (g *SimpleTraceIDGenerator) GenerateSpanID() string {
	return "span-" + generateID()
}

// Span represents a single operation in a trace.
type Span struct {
	TraceID    string
	SpanID     string
	ParentID   string
	Operation  string
	StartTime  time.Time
	EndTime    time.Time
	Attributes map[string]string
	Events     []SpanEvent
	Status     SpanStatus
}

// SpanEvent represents an event within a span.
type SpanEvent struct {
	Time      time.Time
	Name      string
	Attributes map[string]string
}

// SpanStatus represents the status of a span.
type SpanStatus struct {
	Code      string // OK, ERROR, etc.
	Message   string
}

// Tracer creates and manages spans.
type Tracer struct {
	traceIDGenerator TraceIDGenerator
}

// NewTracer creates a new tracer.
func NewTracer(traceIDGenerator TraceIDGenerator) *Tracer {
	if traceIDGenerator == nil {
		traceIDGenerator = NewSimpleTraceIDGenerator()
	}
	return &Tracer{
		traceIDGenerator: traceIDGenerator,
	}
}

// StartSpan starts a new span.
// Attempts to extract trace ID from context first, then generates new if not found.
func (t *Tracer) StartSpan(ctx context.Context, operation string, parentSpanID string) (*Span, context.Context) {
	// Try to get trace ID from context first
	traceID := contextx.TraceID(ctx)
	if traceID == "" && parentSpanID != "" {
		// If no trace ID in context but we have a parent, try to extract from parent span
		if parentSpan := GetSpanFromContext(ctx); parentSpan != nil {
			traceID = parentSpan.TraceID
		}
	}
	if traceID == "" {
		traceID = t.traceIDGenerator.GenerateTraceID()
	}

	spanID := t.traceIDGenerator.GenerateSpanID()
	span := &Span{
		TraceID:    traceID,
		SpanID:     spanID,
		ParentID:   parentSpanID,
		Operation:  operation,
		StartTime:  time.Now(),
		Attributes: make(map[string]string),
		Events:     []SpanEvent{},
		Status:     SpanStatus{Code: "OK"},
	}

	// Store span in context
	ctx = context.WithValue(ctx, spanKey{}, span)
	return span, ctx
}

// EndSpan ends a span.
func (t *Tracer) EndSpan(span *Span) {
	if span == nil {
		return
	}
	span.EndTime = time.Now()
	// In a real implementation, this would send the span to a tracing backend
}

// AddSpanEvent adds an event to a span.
func (t *Tracer) AddSpanEvent(span *Span, name string, attributes map[string]string) {
	if span == nil {
		return
	}
	event := SpanEvent{
		Time:      time.Now(),
		Name:      name,
		Attributes: attributes,
	}
	span.Events = append(span.Events, event)
}

// SetSpanStatus sets the status of a span.
func (t *Tracer) SetSpanStatus(span *Span, code string, message string) {
	if span == nil {
		return
	}
	span.Status.Code = code
	span.Status.Message = message
}

// SetSpanAttribute sets an attribute on a span.
func (t *Tracer) SetSpanAttribute(span *Span, key string, value string) {
	if span == nil {
		return
	}
	if span.Attributes == nil {
		span.Attributes = make(map[string]string)
	}
	span.Attributes[key] = value
}

// spanKey is the context key for storing spans.
type spanKey struct{}

// GetSpanFromContext retrieves the current span from context.
func GetSpanFromContext(ctx context.Context) *Span {
	span, _ := ctx.Value(spanKey{}).(*Span)
	return span
}