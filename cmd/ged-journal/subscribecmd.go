package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/kmpm/ged-journal/internal/misccli"
	"github.com/nats-io/nats.go"
)

type SubscribeCmd struct {
	File SubFileCmd `cmd:"" help:"Save messages to disk"`
}

type SubFileCmd struct {
	Path        string `arg:""  help:"Directory path to save journal files" type:"existingdir"`
	Deflate     bool   `short:"d" help:"Deflate message" default:"false"`
	Nats        string `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string `help:"Nats context" default:""`
	NatsCreds   string `help:"Nats credentials" default:"" type:"existingfile"`
	Subject     string `help:"Subject to save" type:"string" default:">"`
}

func (cmd *SubFileCmd) Run(ctx *clicontext) error {
	slog.Info("Subscribing to journal events", "subject", cmd.Subject, "path", cmd.Path)
	fmt.Println("Subscribing to journal events and saving to", cmd.Path)
	nc, err := misccli.ConnectNATS(cmd.Nats, cmd.NatsContext, cmd.NatsCreds)
	if err != nil {
		return err
	}
	sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
		epoc := time.Now().Unix()
		slog.Debug("Received message", "subject", m.Subject, "epoc", epoc)
		filename := fmt.Sprintf("%d_%s.json", epoc, strings.ReplaceAll(m.Subject, ".", "-"))
		filename = filepath.FromSlash(cmd.Path + "/" + filename)
		data := m.Data
		if cmd.Deflate {
			data, err = compression.Inflate(m.Data)
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
	if err != nil {
		panic(err)
	}
	slog.Info("Waiting for messages")
	fmt.Println("Press Ctrl+C to stop")
	<-waitfor()
	slog.Info("Closing connection")
	err = sub.Unsubscribe()
	if err != nil {
		slog.Warn("Failed to unsubscribe", "error", err)
	} else {
		err = sub.Drain()
		if err != nil {
			slog.Warn("Failed to drain subscription", "error", err)
		}
	}

	nc.Close()
	return nil
}
