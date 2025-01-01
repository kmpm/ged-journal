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
	Save SubSaveCmd `cmd:"" help:"Save journal files to disk"`
}

type SubSaveCmd struct {
	Path        string `arg:""  help:"Directory path to save journal files" type:"existingdir"`
	Subject     string `help:"Subject to save" type:"string" default:">"`
	Nats        string `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string `help:"Nats context" default:""`
}

func (cmd *SubSaveCmd) Run(ctx *clicontext) error {
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
		err := os.WriteFile(filename, m.Data, 0644)
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
