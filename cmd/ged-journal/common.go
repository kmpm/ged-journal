package main

import (
	"log/slog"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	Context string `help:"Nats context name. Leave blank for default"`
	Server  string `help:"Nats server address"`
}

func connect(cfg *Nats) (nc *nats.Conn, err error) {

	if cfg.Server != "" {
		nc, err = nats.Connect(cfg.Server)
	} else {
		nc, err = natscontext.Connect(cfg.Context, nil)
	}
	if err != nil {
		return nil, err
	}
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

	return nc, nil
}
