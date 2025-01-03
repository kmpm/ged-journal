package agent

import (
	"sync"

	"github.com/kmpm/ged-journal/public/messages"
)

type EventHandler func(e messages.Event, fields map[string]interface{}, state *State) (bool, error)

type eventRegistry struct {
	mu       sync.Mutex
	handlers map[string][]EventHandler
}

var (
	registry = &eventRegistry{
		handlers: make(map[string][]EventHandler),
	}
)

func RegisterEventHandlers(handler EventHandler, events ...string) error {
	for _, event := range events {
		if err := registry.RegisterEventHandler(event, handler); err != nil {
			return err
		}
	}
	return nil
}

func RegisterEventHandler(event string, handler EventHandler) error {
	return registry.RegisterEventHandler(event, handler)
}
func (r *eventRegistry) RegisterEventHandler(event string, handler EventHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.handlers[event]; !ok {
		r.handlers[event] = make([]EventHandler, 0)
	}
	r.handlers[event] = append(r.handlers[event], handler)
	return nil
}

// func exampleEventHandler(e messages.Event, fields map[string]any, state *State) error {
// 	// Example event handler logic
// 	return nil
// }
