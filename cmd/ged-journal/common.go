package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"log/slog"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
)

func connect(uri, context string) (nc *nats.Conn, err error) {
	if context != "" {
		nc, err = natscontext.Connect("nats_development", nil)
	} else if uri != "" {
		nc, err = nats.Connect(uri)
	} else {
		return nil, errors.New("no nats server address provided")
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

// Deflate decompresses a zlib compressed byte slice
func deflate(raw []byte) ([]byte, error) {

	r, err := zlib.NewReader(bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	p := make([]byte, 1024)
	out := new(bytes.Buffer)
	n, err := r.Read(p)
	for n > 0 {
		out.Write(p[:n])
		if err != nil {
			break
		}
		n, err = r.Read(p)
	}
	if err != io.EOF {
		return nil, err
	}
	return out.Bytes(), nil

}
