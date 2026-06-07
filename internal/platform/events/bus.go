package events

import (
	"sync"
)

// EventHandler is a function that handles an event.
type EventHandler func(event interface{})

// Bus is a simple in-memory event bus.
type Bus struct {
	mu          sync.RWMutex
	handlers    map[string][]EventHandler
	onceHandlers map[string][]EventHandler
}

// NewBus creates a new event bus.
func NewBus() *Bus {
	return &Bus{
		handlers:     make(map[string][]EventHandler),
		onceHandlers: make(map[string][]EventHandler),
	}
}

// Publish publishes an event to all registered handlers.
func (b *Bus) Publish(eventType string, event interface{}) {
	b.mu.RLock()
	handlers := b.handlers[eventType]
	onceHandlers := b.onceHandlers[eventType]
	b.mu.RUnlock()

	// Call regular handlers
	for _, handler := range handlers {
		handler(event)
	}

	// Call once handlers and then remove them
	for _, handler := range onceHandlers {
		handler(event)
	}
	b.mu.Lock()
	delete(b.onceHandlers, eventType)
	b.mu.Unlock()
}

// Subscribe registers a handler for the given event type.
func (b *Bus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// SubscribeOnce registers a handler that will be called only once for the given event type.
func (b *Bus) SubscribeOnce(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.onceHandlers[eventType] = append(b.onceHandlers[eventType], handler)
}

// Unsubscribe removes a handler for the given event type.
// Note: This implementation does not actually remove the handler because we don't have a way to identify it.
// In a real implementation, you might return a subscription object that can be used to unsubscribe.
func (b *Bus) Unsubscribe(eventType string, handler EventHandler) {
	// TODO: implement if needed
}