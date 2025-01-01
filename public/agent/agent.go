package agent

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/nats-io/nats.go"
)

type statusSub struct {
	Chan chan *Status
	Sub  *nats.Subscription
}

type dataSub struct {
	Chan chan []byte
	Sub  *nats.Subscription
}

type Agent struct {
	open           bool
	nc             *nats.Conn
	statusSubs     []statusSub
	dataSubs       []dataSub
	mu             sync.Mutex
	currentStatus  *Status
	prefix         string
	currentSystem  string
	currentStation string
}

func New(nc *nats.Conn, prefix string) (a *Agent, err error) {
	if nc == nil {
		return nil, errors.New("no nats connection provided")
	}
	a = &Agent{
		nc:         nc,
		statusSubs: make([]statusSub, 0),
		dataSubs:   make([]dataSub, 0),
		open:       true,
	}
	return a, nil
}

func (a *Agent) Close() {
	//TODO: close all subscriptions
	if !a.open {
		return
	}
	a.mu.Lock()
	for _, s := range a.statusSubs {
		s.Sub.Unsubscribe()
		close(s.Chan)
	}
	a.statusSubs = nil
	for _, s := range a.dataSubs {
		s.Sub.Unsubscribe()
		close(s.Chan)
	}
	a.dataSubs = nil
	a.open = false
	a.mu.Unlock()
}

func (a *Agent) Status() chan *Status {
	if !a.open {
		panic("Agent is closed")
	}
	ch := make(chan *Status)
	s, err := a.nc.Subscribe(a.prefix+"status", func(msg *nats.Msg) {
		var (
			stat *Status
			err  error
			data []byte
		)
		if msg.Header.Get("Encoding") == "zlib" {
			data, err = compression.Inflate(msg.Data)
			if err != nil {
				slog.Error("Failed to decompress message", "error", err)
				return
			}
		} else {
			data = msg.Data
		}

		stat, err = GetStatusFromBytes(data)
		if err != nil {
			slog.Error("Failed to decode message", "error", err)
			return
		}
		a.mu.Lock()
		a.currentStatus = stat
		a.mu.Unlock()
		ch <- stat
	})
	if err != nil {
		close(ch)
	}
	a.mu.Lock()
	a.statusSubs = append(a.statusSubs, statusSub{Chan: ch, Sub: s})
	a.mu.Unlock()
	return ch
}

func (a *Agent) Message(subject string) chan []byte {
	if !a.open {
		panic("Agent is closed")
	}
	// TODO: check valid name
	ch := make(chan []byte)
	s, err := a.nc.Subscribe(a.prefix+subject, onMessageHandler(ch))
	if err != nil {
		close(ch)
	}
	a.mu.Lock()
	a.dataSubs = append(a.dataSubs, dataSub{Chan: ch, Sub: s})
	a.mu.Unlock()
	return ch

}

func onMessageHandler(ch chan []byte) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		var (
			data []byte
			err  error
		)
		if msg.Header.Get("Encoding") == "zlib" {
			data, err = compression.Inflate(msg.Data)
			if err != nil {
				slog.Error("Failed to decompress message", "error", err)
				return
			}
		} else {
			data = msg.Data
		}
		ch <- data
	}
}
