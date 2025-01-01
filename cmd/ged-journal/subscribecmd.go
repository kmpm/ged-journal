package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type SubscribeCmd struct {
	File SubFileCmd `cmd:"" help:"Save messages to disk"`
}

type SubFileCmd struct {
	Path        string `arg:""  help:"Directory path to save journal files" type:"existingdir"`
	Subject     string `help:"Subject to save" type:"string" default:">"`
	Nats        string `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string `help:"Nats context" default:""`
	Deflate     bool   `short:"d" help:"Deflate message" default:"false"`
}

func (cmd *SubFileCmd) Run(ctx *clicontext) error {
	slog.Info("Subscribing to journal events", "subject", cmd.Subject, "path", cmd.Path)
	nc, err := connect(cmd.Nats, cmd.NatsContext)
	if err != nil {
		return err
	}
	nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
		epoc := time.Now().Unix()
		slog.Debug("Received message", "subject", m.Subject, "epoc", epoc)
		filename := fmt.Sprintf("%s_%d.json", strings.ReplaceAll(m.Subject, ".", "-"), epoc)
		filename = filepath.FromSlash(cmd.Path + "/" + filename)
		data := m.Data
		if cmd.Deflate {
			data, err = deflate(m.Data)
			if err != nil {
				slog.Error("Failed to deflate message", "error", err)
				return
			}
		}
		err := os.WriteFile(filename, data, 0644)
		if err != nil {
			slog.Error("Failed to write file", "file", filename, "error", err)
		}
	})
	slog.Info("Waiting for messages")
	<-waitfor()
	slog.Info("Closing connection")
	nc.Close()
	return nil
}
