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
	log.Printf("api runtime bootstrapped: env=%s addr=%s", deps.Config.Env, deps.Config.API.HTTPAddress)

	// Sprint 1 skeleton only:
	// - wire logger, database, redis, migrations, health endpoints
	// - add gin router and middleware chain under /api/v1
}
