package app

import (
	"context"
	"log"

	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/observability"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/events"
)

// RuntimeDeps holds all runtime dependencies.
type RuntimeDeps struct {
	Config    config.Config
	Logger    *logger.Logger
	Postgres  *database.Postgres
	Redis     *database.Redis
	EventBus  events.EventBus
	Metrics   observability.Metrics
	Tracer    *observability.Tracer
	Shutdown  func(context.Context) error
}

// BuildRuntimeDeps creates all runtime dependencies.
func BuildRuntimeDeps(cfg config.Config) RuntimeDeps {
	// Initialize logger
	var err error
	LoggerInstance, err := logger.New(cfg.Observability.LogLevel, cfg.Observability.ServiceName)
	if err != nil {
		// Use standard logger as fallback since our logger isn't initialized yet
		log.Fatalf("failed to initialize logger: %v", err)
	}
	log := LoggerInstance

	// Initialize PostgreSQL
	pg, err := database.NewPostgres(cfg)
	if err != nil {
		log.Sugar().Fatalf("failed to connect to postgres: %v", err)
	}

	// Initialize Redis
	redis, err := database.NewRedis(cfg)
	if err != nil {
		log.Sugar().Fatalf("failed to connect to redis: %v", err)
	}

	// Initialize event bus
	eventBus := events.NewInMemoryBus()

	// Initialize metrics collector
	metricsCollector := observability.NewMetricsCollector()
	metricsAdapter := observability.NewMetricsAdapter(metricsCollector)

	// Initialize tracer
	traceIDGenerator := observability.NewSimpleTraceIDGenerator()
	tracer := observability.NewTracer(traceIDGenerator)

	// Shutdown function
	shutdown := func(ctx context.Context) error {
		log.Sugar().Info("shutting down services...")

		// Close database connections
		if err := pg.Close(); err != nil {
			log.Sugar().Errorf("error closing postgres connection: %v", err)
		}
		if err := redis.Close(); err != nil {
			log.Sugar().Errorf("error closing redis connection: %v", err)
		}

		// Additional shutdown logic can be added here
		return nil
	}

	log.Sugar().Infof("runtime dependencies initialized: env=%s", cfg.Env)

	return RuntimeDeps{
		Config:   cfg,
		Logger:   log,
		Postgres: pg,
		Redis:    redis,
		EventBus: eventBus,
		Metrics:  metricsAdapter,
		Tracer:   tracer,
		Shutdown: shutdown,
	}
}