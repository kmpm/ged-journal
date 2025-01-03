package agent

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/kmpm/ged-journal/public/messages"
	"github.com/nats-io/nats.go"
)

type State struct {
	changed   bool
	Timestamp string
	Docked    bool
	Game      struct {
		Language string
		Version  string
		Build    string
		Horizons bool
		Odyssey  bool
	}
	Cmdr     string
	FID      string
	System   messages.System
	Station  messages.Station
	Body     messages.Body
	GameMode string
	Credits  int
	Status   *messages.StatusEvent
	Ship     struct {
		Name string
		ID   int
	}
}

func (i *State) Flags() (f messages.StatusFlags) {
	messages.ExpandFlags(i.Status.RawFlags, &f)
	return
}

type dataSub struct {
	Subject string
	Sub     *nats.Subscription
}

type Agent struct {
	isOpen   bool
	nc       *nats.Conn
	dataSubs []dataSub
	mu       sync.Mutex

	collectPrefix string
	agentPrefix   string

	state *State
	tick  *time.Ticker
}

func New(nc *nats.Conn, collectPrefix, agentPrefix string) (a *Agent, err error) {
	if nc == nil {
		return nil, errors.New("no nats connection provided")
	}
	a = &Agent{
		nc:            nc,
		dataSubs:      make([]dataSub, 0),
		collectPrefix: collectPrefix,
		agentPrefix:   agentPrefix,
		state: &State{
			Status: &messages.StatusEvent{},
		},
	}
	go a.open()
	return a, nil
}

func (a *Agent) open() {
	a.mu.Lock()
	a.isOpen = true
	a.mu.Unlock()

	a.tick = time.NewTicker(2 * time.Second)
	go func() {
		for range a.tick.C {
			fmt.Println("--------------")
			fmt.Println("Timestamp:", a.state.Timestamp)
			fmt.Println("Cmdr:", a.state.Cmdr)
			fmt.Println("FID:", a.state.FID)
			fmt.Printf("Game: %+v\n", a.state.Game)
			fmt.Println("Docked:", a.state.Docked)
			fmt.Println("Flags", a.state.Status.RawFlags)
			fmt.Printf("System: %+v\n", a.state.System)
			fmt.Printf("Body: %+v\n", a.state.Body)
			fmt.Printf("Station: %+v\n", a.state.Station)
			fmt.Println("")
		}
	}()

	slog.Info("Checking for status messages")

	a.Message(">", func(data []byte) {
		// slog.Info("Received event", "data", string(data))
		evt := messages.Event{}
		fields, err := messages.GetEventComponent(data, &evt)
		if err != nil {
			slog.Error("Failed to decode event message", "error", err)
			return
		}

		// loop through all the eventhandlers in registry
		registry.mu.Lock()
		defer registry.mu.Unlock()
		a.mu.Lock()
		defer a.mu.Unlock()
		count := 0
		for _, handler := range registry.handlers[evt.Event] {
			_, err := handler(evt, fields, a.state)
			if err != nil {
				slog.Error("Failed to handle event", "event", evt.Event, "error", err)
			}
			a.state.Timestamp = evt.Timestamp
			count++
			//todo: check if state has changed
		}
		// if count == 0 {
		// 	slog.Debug("No handlers registered for event", "event", evt.Event)
		// }
	})
}

func (a *Agent) Close() {
	//TODO: close all subscriptions
	a.tick.Stop()
	if !a.isOpen {
		return
	}
	a.mu.Lock()

	for _, s := range a.dataSubs {
		slog.Debug("Unsubscribing from message", "subject", s.Subject)
		s.Sub.Unsubscribe()
	}
	a.dataSubs = nil
	a.isOpen = false
	a.mu.Unlock()
	slog.Info("Agent closed", "state", a.state)
}

type MessageHandler func([]byte)

func (a *Agent) Message(subject string, cb MessageHandler) error {
	if !a.isOpen {
		return errors.New("Agent is closed")
	}
	// TODO: check valid subjects
	subject = a.collectPrefix + subject
	slog.Debug("Subscribing to message", "subject", subject)
	s, err := a.nc.Subscribe(subject, onMessageHandler(cb))
	if err != nil {
		slog.Info("Failed to subscribe to message", "subject", subject, "error", err)
		return err
	}
	a.mu.Lock()
	a.dataSubs = append(a.dataSubs, dataSub{Subject: subject, Sub: s})
	a.mu.Unlock()
	return nil
}

func onMessageHandler(cb MessageHandler) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		var (
			data []byte
			err  error
		)
		if msg.Header.Get("Content-Encoding") == "zlib" {
			data, err = compression.Inflate(msg.Data)
			if err != nil {
				slog.Error("Failed to decompress message", "error", err)
				return
			}
		} else {
			data = msg.Data
		}
		cb(data)
	}
}
