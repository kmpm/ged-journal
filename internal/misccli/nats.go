package misccli

import (
	"log/slog"
	"path/filepath"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
)

type NatsFactory func() (*nats.Conn, error)

func ConnectNATS(uri, context, creds string) (*nats.Conn, error) {
	var nc *nats.Conn
	var err error
	if context != "" {
		slog.Info("Using nats context", "context", context)
		nc, err = Connect(UsingContextg(context))
	} else {
		if creds != "" {
			slog.Info("Using nats creds", "creds", creds, "uri", uri)
			nc, err = Connect(UsingCreds(uri, creds))
		} else {
			slog.Info("Using nats uri", "uri", uri)
			nc, err = Connect(UsingURI(uri))
		}
	}
	if err == nil {
		slog.Info("Connected to nats server", "servers", nc.Servers())
	} else {
		slog.Error("Failed to connect to nats server", "error", err)
	}
	return nc, err
}

func UsingContextg(uri string) NatsFactory {
	return func() (*nats.Conn, error) {
		return natscontext.Connect("nats_development", nil)
	}
}

func UsingURI(uri string) NatsFactory {
	return func() (*nats.Conn, error) {
		return nats.Connect(uri)
	}
}

func UsingCreds(uri string, creds string) NatsFactory {
	return func() (*nats.Conn, error) {
		return nats.Connect(uri, nats.UserCredentials(filepath.FromSlash(creds)))
	}
}

func Connect(nf NatsFactory) (nc *nats.Conn, err error) {
	nc, err = nf()
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
