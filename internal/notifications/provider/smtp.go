package provider

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/AARCSX/AARCSX_Forge/internal/notifications"
)

// SMTPProvider implements the Provider interface for SMTP email sending.
type SMTPProvider struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewSMTPProvider creates a new SMTP provider instance.
func NewSMTPProvider(host string, port int, username, password, from string) (*SMTPProvider, error) {
	if host == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}
	if port <= 0 {
		return nil, fmt.Errorf("SMTP port must be positive")
	}
	if from == "" {
		return nil, fmt.Errorf("SMTP from address is required")
	}

	return &SMTPProvider{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}, nil
}

// SendEmail sends an email via SMTP.
func (p *SMTPProvider) SendEmail(ctx context.Context, in notifications.EmailMessage) error {
	// Validate inputs
	if in.To == "" {
		return fmt.Errorf("email recipient is required")
	}
	if in.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if in.Body == "" {
		return fmt.Errorf("email body is required")
	}

	// Set up authentication
	auth := smtp.PlainAuth("", p.username, p.password, p.host)

	// Create message
	msg := []byte(
		"From: " + p.from + "\r\n" +
			"To: " + in.To + "\r\n" +
			"Subject: " + in.Subject + "\r\n" +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
			in.Body + "\r\n",
	)

	// Send email
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	if err := smtp.SendMail(addr, auth, p.from, []string{in.To}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}