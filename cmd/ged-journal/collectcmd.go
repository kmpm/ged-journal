package main

import (
	"fmt"
	"log/slog"

	"github.com/kmpm/ged-journal/internal/compression"
	"github.com/kmpm/ged-journal/internal/misccli"
	"github.com/kmpm/ged-journal/public/collector"
	"github.com/nats-io/nats.go"
)

type CollectCmd struct {
	BasePath    string `arg:"" help:"Path to application log files" default:"${basepath}"`
	Nats        string `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string `help:"Nats context" default:""`
	NatsCreds   string `help:"Nats credentials" default:""`
}

func (cmd *CollectCmd) Run(cc *clicontext) error {
	slog.Info("Running Collect")
	nc, err := misccli.ConnectNATS(cmd.Nats, cmd.NatsContext, cmd.NatsCreds)
	if err != nil {
		return err
	}

	pub := func(subject string, data []byte, compress bool) (err error) {
		msg := nats.NewMsg("ged.collector." + subject)
		if compress {
			data, err = compression.Deflate(data)
			if err != nil {
				panic(err)
			}
			msg.Header.Set("Content-Encoding", "zlib")
		}
		msg.Data = data
		return nc.PublishMsg(msg)
	}

	a, err := collector.New(cmd.BasePath, pub)
	if err != nil {
		return err
	}
	defer a.Close()
	fmt.Println("Collecting data")
	fmt.Println("Press Ctrl+C to stop")
	stop := waitfor()
	<-stop
	slog.Info("Shutting down ged-journal")
	return nil
}
