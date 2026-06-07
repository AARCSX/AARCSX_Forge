package observability

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]any)
	Error(ctx context.Context, msg string, err error, fields map[string]any)
}

type AuditLogger interface {
	Write(ctx context.Context, event AuditEvent) error
}

type Metrics interface {
	IncCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
}

type AuditEvent struct {
	Action   string
	ActorID  string
	TenantID string
	TraceID  string
	Resource string
	Status   string
}
