package notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/AARCSX/AARCSX_Forge/internal/logger"
)

// NotificationService implements the Service interface for notification operations.
type NotificationService struct {
	provider Provider
	repo     Repository
	logger   *logger.Logger
}

// NewNotificationService creates a new notification service instance.
func NewNotificationService(logger *logger.Logger) (*NotificationService, error) {
	// For now, we'll return an error since we need to initialize the provider based on config
	// In a real implementation, this would take config and initialize the appropriate provider
	return nil, errors.New("notification service requires provider initialization")
}

// WithProvider sets the provider for the notification service.
func (s *NotificationService) WithProvider(p Provider) *NotificationService {
	s.provider = p
	return s
}

// WithRepository sets the repository for the notification service.
func (s *NotificationService) WithRepository(r Repository) *NotificationService {
	s.repo = r
	return s
}

// QueueEmail queues an email for sending and creates a delivery job record.
func (s *NotificationService) QueueEmail(ctx context.Context, in QueueEmailInput) (DeliveryJob, error) {
	// Validate input
	if in.TenantID == "" {
		return DeliveryJob{}, errors.New("tenant ID is required")
	}
	if in.To == "" {
		return DeliveryJob{}, errors.New("recipient email is required")
	}
	if in.Subject == "" {
		return DeliveryJob{}, errors.New("email subject is required")
	}
	if in.Body == "" {
		return DeliveryJob{}, errors.New("email body is required")
	}

	// Create delivery job record
	deliveryJob := DeliveryJob{
		TenantID: in.TenantID,
		Channel:  "email",
		Status:   "queued",
	}

	// Save to repository if available
	if s.repo != nil {
		var err error
		deliveryJob, err = s.repo.CreateDelivery(ctx, deliveryJob)
		if err != nil {
			s.logger.Sugar().Errorw("failed to create delivery job record", "error", err, "tenant_id", in.TenantID)
			return DeliveryJob{}, fmt.Errorf("failed to create delivery job: %w", err)
		}
	}

	// Attempt to send email via provider
	if s.provider != nil {
		emailMsg := EmailMessage{
			To:      in.To,
			Subject: in.Subject,
			Body:    in.Body,
		}

		if err := s.provider.SendEmail(ctx, emailMsg); err != nil {
			// Update job status to failed if we have a repository
			if s.repo != nil && deliveryJob.ID != "" {
				updateErr := s.repo.MarkFailed(ctx, deliveryJob.ID, err.Error())
				if updateErr != nil {
					s.logger.Sugar().Warnw("failed to mark delivery job as failed", "error", updateErr, "job_id", deliveryJob.ID)
				}
			}
			s.logger.Sugar().Errorw("failed to send email", "error", err, "tenant_id", in.TenantID, "to", in.To)
			return DeliveryJob{}, fmt.Errorf("failed to send email: %w", err)
		}

		// Mark as delivered if we have a repository
		if s.repo != nil && deliveryJob.ID != "" {
			if err := s.repo.MarkDelivered(ctx, deliveryJob.ID); err != nil {
				s.logger.Sugar().Warnw("failed to mark delivery job as delivered", "error", err, "job_id", deliveryJob.ID)
				// Don't return error here as the email was sent successfully
			}
			deliveryJob.Status = "delivered"
		}
	}

	return deliveryJob, nil
}