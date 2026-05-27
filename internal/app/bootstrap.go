package app

import (
	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/events"
)

type RuntimeDeps struct {
	Config   config.Config
	EventBus events.Bus
}

func BuildRuntimeDeps(cfg config.Config) RuntimeDeps {
	return RuntimeDeps{
		Config:   cfg,
		EventBus: events.NewInMemoryBus(),
	}
}
