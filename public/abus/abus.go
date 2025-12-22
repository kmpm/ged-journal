// Package abus provides a simple wrapper around the NATS messaging system.
package abus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/kmpm/ged-journal/internal/metrics"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	EncodingHeader = "content-encoding"
)

type Config struct {
	Server  string `help:"NATS server URL" env:"SERVER"`                        // NATS server URL
	Context string `help:"NATS context. Leave empty for default" env:"CONTEXT"` // NATS context
	User    string `help:"NATS user" env:"USER"`                                // NATS user
	Pass    string `help:"NATS password" env:"PASS"`                            // NATS password
	Token   string `help:"NATS authentication token" env:"TOKEN"`               // NATS authentication token
	Nkey    string `help:"NATS NKey file path" type:"file" env:"NKEY"`          // NATS NKey
}

func (c Config) String() string {
	return fmt.Sprintf("Server: %s, User: %v, Pass: %v, Token: %v, Nkey: %s", c.Server, c.User != "", c.Pass != "", c.Token != "", c.Nkey)
}

type ABus struct {
	nc           *nats.Conn
	debugEnabled bool
}

func Connect(c Config) (*ABus, error) {
	opts := []nats.Option{
		nats.Name("GED-Trade"),
		nats.Timeout(5 * time.Second),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(1 * time.Second),
	}
	if c.User != "" && c.Pass != "" {
		opts = append(opts, nats.UserInfo(c.User, c.Pass))
	}
	if c.Token != "" {
		opts = append(opts, nats.Token(c.Token))
	}
	if c.Nkey != "" {
		nkey, err := nats.NkeyOptionFromSeed(c.Nkey)
		if err != nil {
			return nil, err
		}
		opts = append(opts, nkey)
	}
	var nc *nats.Conn
	var err error
	if c.Server != "" {
		nc, err = nats.Connect(c.Server, opts...)
	} else {
		nc, err = natscontext.Connect(c.Context, opts...)
	}
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	//TODO: Add logging, error handling, and callbacks.
	nc.SetClosedHandler(func(_ *nats.Conn) {
		slog.Error("nats connection closed")
	})
	nc.SetDisconnectHandler(func(_ *nats.Conn) {
		slog.Error("nats connection disconnected")
	})
	nc.SetDisconnectErrHandler(func(_ *nats.Conn, err error) {
		slog.Error("nats connection disconnected", "error", err)
	})
	nc.SetReconnectHandler(func(_ *nats.Conn) {
		slog.Info("nats connection reconnected")
	})

	return &ABus{
		nc:           nc,
		debugEnabled: slog.Default().Enabled(context.Background(), slog.LevelDebug),
	}, nil
}

func (a *ABus) Conn() *nats.Conn {
	return a.nc
}

func (a *ABus) Close() error {
	err := a.nc.Drain()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe subscribes to a subject with optional message inflation.
// If inflate is true or the message header contains "Encoding: zlib", the message data will be inflated using zlib.
func (a *ABus) Subscribe(subject string, inflate bool, callback func(msg *nats.Msg)) (*nats.Subscription, error) {
	metInboundSizeSaved := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ged_subscribe_size_saved",
		Help: "The amount of bytes because of message compression on inbound messages",
	})
	metrics.Registry.MustRegister(metInboundSizeSaved)
	sub, err := a.nc.Subscribe(subject, func(msg *nats.Msg) {
		sizeBefore := len(msg.Data)
		if inflate || msg.Header.Get(EncodingHeader) == "zlib" {
			if a.debugEnabled {
				slog.Debug("inflating message", "subject", msg.Subject, "encoding", msg.Header.Get(EncodingHeader))
			}
			data, err := compression.Inflate(msg.Data)
			if err != nil {
				slog.Error("failed to inflate message", "error", err)
				return
			}
			sizeDiff := len(data) - sizeBefore
			metInboundSizeSaved.Add(float64(sizeDiff))
			if a.debugEnabled {
				slog.Debug("inflated message", "subject", msg.Subject, "size_before", sizeBefore, "size_after", len(data), "size_diff", sizeDiff)
			}
			msg.Header.Del(EncodingHeader)
			msg.Data = data
		}
		callback(msg)
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (a *ABus) QueueSubscribe(subject string, inflate bool, queue string, callback func(msg *nats.Msg)) (*nats.Subscription, error) {
	metInboundSizeSaved := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ged_queue_subscribe_size_saved",
		Help: "The amount of bytes because of message compression on inbound messages",
	})
	metrics.Registry.MustRegister(metInboundSizeSaved)
	consumer, err := a.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		sizeBefore := len(msg.Data)
		if inflate || msg.Header.Get(EncodingHeader) == "zlib" {
			if a.debugEnabled {
				slog.Debug("inflating message", "subject", msg.Subject, "encoding", msg.Header.Get(EncodingHeader))
			}
			data, err := compression.Inflate(msg.Data)
			if err != nil {
				slog.Error("failed to inflate message", "error", err)
				return
			}
			sizeAfter := len(data)
			sizeDiff := sizeBefore - sizeAfter
			metInboundSizeSaved.Add(float64(sizeDiff))
			if a.debugEnabled {
				slog.Debug("inflated message", "subject", msg.Subject, "size_before", sizeBefore, "size_after", sizeAfter, "size_diff", sizeDiff)
			}
			msg.Header.Del(EncodingHeader)
			msg.Data = data
		}
		callback(msg)
	})
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (a *ABus) Publish(subject string, data []byte, deflate bool) error {
	var err error
	msg := nats.NewMsg(subject)
	if deflate {
		if a.debugEnabled {
			slog.Debug("deflating message", "subject", subject)
		}
		data, err = compression.Deflate(data)
		if err != nil {
			return err
		}
		msg.Header.Set(EncodingHeader, "zlib")
	}
	msg.Data = data
	return a.nc.PublishMsg(msg)
}

func (a *ABus) PublishCompressed(subject string, data []byte) error {
	return a.Publish(subject, data, true)
}
