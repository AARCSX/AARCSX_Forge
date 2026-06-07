package events

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Event interface {
	Name() string
	Version() int
	OccurredAt() time.Time
}

type Handler func(context.Context, Event) error

type EventBus interface {
	Subscribe(eventName string, handler Handler)
	Publish(ctx context.Context, event Event) error
}

type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *InMemoryBus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := append([]Handler{}, b.handlers[event.Name()]...)
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return fmt.Errorf("handle event %s: %w", event.Name(), err)
		}
	}
	return nil
}
