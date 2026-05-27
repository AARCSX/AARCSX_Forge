package main

import (
	"log"

	"github.com/AARCSX/AARCSX_Forge/internal/app"
	"github.com/AARCSX/AARCSX_Forge/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	deps := app.BuildRuntimeDeps(cfg)
	log.Printf("worker runtime bootstrapped: env=%s queue_concurrency=%d", deps.Config.Env, deps.Config.Queue.Concurrency)

	// Sprint 4 skeleton only:
	// - wire Asynq server, queue handlers, retries, and dead-letter policies
}
