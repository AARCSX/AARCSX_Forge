package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/AARCSX/AARCSX_Forge/internal/app"
	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/AARCSX/AARCSX_Forge/internal/notifications"
	"github.com/AARCSX/AARCSX_Forge/internal/notifications/provider"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
)

// EmailTaskPayload defines the payload for email processing tasks.
type EmailTaskPayload struct {
	TenantID string `json:"tenant_id"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

// handleEmailTask processes email sending tasks.
func handleEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	// Get logger from context (or create a fallback)
	var log *logger.Logger
	log = &logger.Logger{} // fallback
	if loggerFromCtx := ctx.Value("logger"); loggerFromCtx != nil {
		if l, ok := loggerFromCtx.(*logger.Logger); ok {
			log = l
		}
	}

	// Initialize SMTP provider (in practice, this would come from config/di)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")

	smtpPort := 587 // default
	if smtpPortStr != "" {
		var err error
		smtpPort, err = strconv.Atoi(smtpPortStr)
		if err != nil {
			smtpPort = 587
		}
	}

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || smtpFrom == "" {
		log.Sugar().Warnw("SMTP configuration incomplete, skipping email", "host", smtpHost, "user", smtpUser, "from", smtpFrom)
		return nil // Don't fail the task if SMTP not configured
	}

	smtpProvider, err := provider.NewSMTPProvider(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom)
	if err != nil {
		log.Sugar().Errorw("Failed to initialize SMTP provider", "error", err)
		return err
	}

	// Send email
	emailMsg := notifications.EmailMessage{
		To:      p.To,
		Subject: p.Subject,
		Body:    p.Body,
	}

	if err := smtpProvider.SendEmail(ctx, emailMsg); err != nil {
		log.Sugar().Errorw("Failed to send email", "error", err, "to", p.To)
		return err
	}

	log.Sugar().Infow("Email sent successfully", "to", p.To)
	return nil
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Build runtime dependencies (for logger, etc.)
	deps := app.BuildRuntimeDeps(cfg)

	// Configure Asynq scheduler
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.URL},
		asynq.Config{
			Concurrency: int(cfg.Queue.Concurrency),
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				deps.Logger.Sugar().Errorw("task failed", "error", err, "task_type", task.Type())
			}),
		},
	)

	// Handle shutdown signals
	mux := asynq.NewServeMux()
	mux.HandleFunc("email:send", handleEmailTask)

	// Run worker
	go func() {
		if err := srv.Start(mux); err != nil {
			log.Fatalf("failed to start Asynq server: %v", err)
		}
	}()

	log.Printf("Worker started: env=%s, redis=%s, concurrency=%d",
		cfg.Env, cfg.Redis.URL, cfg.Queue.Concurrency)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down worker...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown()

	// Shutdown dependencies
	if err := deps.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down dependencies: %v", err)
	}

	log.Println("Worker stopped")
}