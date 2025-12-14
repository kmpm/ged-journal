package main

import (
	"log/slog"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/kmpm/ged-journal/public/collector"
	"github.com/nats-io/nats.go"
)

type CollectCmd struct {
	BasePath string `arg:"" help:"Path to application log files" default:"${basepath}"`
	Nats     Nats   `help:"Nats server address" embed:"" prefix:"nats."`
}

func (cmd *CollectCmd) Run(cc *clicontext) error {
	slog.Info("Running Collect")
	nc, err := connect(&cmd.Nats)
	if err != nil {
		return err
	}

	pub := func(subject string, data []byte, compress bool) (err error) {
		msg := nats.NewMsg("ged." + subject)
		if compress {
			data, err = compression.Deflate(data)
			if err != nil {
				panic(err)
			}
			msg.Header.Set("Encoding", "zlib")
		}
		msg.Data = data
		return nc.PublishMsg(msg)
	}

	a, err := collector.New(cmd.BasePath, pub)
	if err != nil {
		return err
	}
	defer a.Close()

	stop := waitfor()
	<-stop
	slog.Info("Shutting down ged-journal")
	return nil
}
