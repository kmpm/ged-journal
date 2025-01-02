package agent

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/kmpm/ged-journal/public/messages"
	"github.com/nats-io/nats.go"
)

type dataSub struct {
	Subject string
	Sub     *nats.Subscription
}

type Agent struct {
	isOpen         bool
	nc             *nats.Conn
	dataSubs       []dataSub
	mu             sync.Mutex
	currentStatus  *messages.Status
	collectPrefix  string
	agentPrefix    string
	currentSystem  string
	currentStation string
	currentHeader  messages.FileHeader
	currentCmdr    messages.Commander
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
	}
	go a.open()
	return a, nil
}

func (a *Agent) open() {
	a.mu.Lock()
	a.isOpen = true
	a.mu.Unlock()
	slog.Info("Checking for status messages")
	a.Status(func(stat messages.Status) {
		a.mu.Lock()
		a.currentStatus = &stat
		a.mu.Unlock()
		slog.Info("Received status", "status", stat)
	})

	slog.Info("Checking for fileheader messages")
	a.Message("journal.event.fileheader", func(data []byte) {
		//TODO: validate fileheader
		header := messages.FileHeader{}
		err := messages.GetFileHeader(data, &header)
		if err != nil {
			slog.Error("Failed to decode fileheader message", "error", err)
			return
		}
		a.mu.Lock()
		a.currentHeader = header
		a.mu.Unlock()
		slog.Info("Received fileheader", "header", header)
	})
	slog.Info("Checking for commander messages")
	a.Message("journal.event.commander", func(data []byte) {
		cmdr := messages.Commander{}
		err := messages.GetCommander(data, &cmdr)
		if err != nil {
			slog.Error("Failed to decode commander message", "error", err)
			return
		}
		a.mu.Lock()
		a.currentCmdr = cmdr
		a.mu.Unlock()
		slog.Info("Received commander", "commander", cmdr)
	})
}

func (a *Agent) Close() {
	//TODO: close all subscriptions
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
}

type StatusHandler func(messages.Status)
type MessageHandler func([]byte)

func (a *Agent) Status(cb StatusHandler) error {
	if !a.isOpen {
		return errors.New("Agent is closed")
	}
	status := messages.Status{}
	a.Message("global.status", func(data []byte) {
		err := messages.GetStatus(data, &status)
		if err != nil {
			slog.Error("Failed to decode status message", "error", err)
			return
		}
		cb(status)
	})
	return nil
}

func (a *Agent) Message(subject string, cb MessageHandler) error {
	if !a.isOpen {
		return errors.New("Agent is closed")
	}
	// TODO: check valid subjects
	subject = a.collectPrefix + subject
	slog.Debug("Subscribing to message", "subject", subject)
	s, err := a.nc.Subscribe(subject, onMessageHandler(cb))
	if err != nil {
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
