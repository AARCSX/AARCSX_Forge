package notifications

import "context"

type Service interface {
	QueueEmail(ctx context.Context, in QueueEmailInput) (DeliveryJob, error)
}

type Provider interface {
	SendEmail(ctx context.Context, in EmailMessage) error
}

type Repository interface {
	CreateDelivery(ctx context.Context, in DeliveryJob) (DeliveryJob, error)
	MarkDelivered(ctx context.Context, jobID string) error
	MarkFailed(ctx context.Context, jobID string, reason string) error
}

type QueueEmailInput struct {
	TenantID string
	To       string
	Subject  string
	Body     string
}

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

type DeliveryJob struct {
	ID       string
	TenantID string
	Channel  string
	Status   string
}
