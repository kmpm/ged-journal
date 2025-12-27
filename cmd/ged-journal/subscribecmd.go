package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/kmpm/ged-common/pkg/nbus"
	"github.com/kmpm/ged-journal/internal/metrics"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
)

type SubscribeCmd struct {
	File SubFileCmd `cmd:"" help:"Save messages to disk"`
}

type SubFileCmd struct {
	Path    string      `arg:""  help:"Directory path to save journal files" type:"existingdir"`
	Subject string      `help:"Subject to save" short:"s" type:"string" default:"ged-journal.>"`
	Nats    nbus.Config `help:"Nats Configuration" prefix:"nats." embed:"" envprefix:"NATS_"`
	Inflate bool        `short:"i" help:"Force inflation of message using zlib" default:"false"`
}

func (cmd *SubFileCmd) Run(ctx *clicontext) error {
	slog.Info("Subscribing to journal events", "subject", cmd.Subject, "path", cmd.Path)

	nb, err := nbus.Connect(cmd.Nats)
	if err != nil {
		return err
	}
	nb.Subscribe(cmd.Subject, cmd.Inflate, saveAsFileMiddleware(cmd.Path))
	slog.Info("Waiting for messages")
	<-waitfor()
	slog.Info("Closing connection")
	nb.Close()
	return nil
}

func saveAsFileMiddleware(path string) func(*nats.Msg) {
	metSaved := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ged_messages_saved",
		Help: "Number of messages saved to file",
	}, []string{"subject"})
	metrics.Registry.MustRegister(metSaved)

	return func(msg *nats.Msg) {
		epoc := time.Now().Unix()
		filename := fmt.Sprintf("%s_%d.json", strings.ReplaceAll(msg.Subject, ".", "-"), epoc)
		filename = filepath.FromSlash(path + "/" + filename)
		slog.Debug("Received message", "subject", msg.Subject, "epoc", epoc, "filename", filename)
		err := os.WriteFile(filename, msg.Data, 0644)
		if err != nil {
			slog.Error("Failed to write file", "file", filename, "error", err)
		}
		metSaved.WithLabelValues(msg.Subject).Inc()
	}
}
