package main

import (
	"log/slog"

	"codeberg.org/kmpm/ged-common/pkg/nbus"
	"github.com/kmpm/ged-journal/public/collector"
)

type CollectCmd struct {
	BasePath      string      `arg:"" help:"Path to application log files" default:"${basepath}"`
	Nats          nbus.Config `embed:"" prefix:"nats." envprefix:"NATS_"`
	SubjectPrefix string      `help:"Prefix for journal event subjects" default:"ged-journal."`
}

func (cmd *CollectCmd) Run(cc *clicontext) error {
	slog.Info("Running Collect")
	nb, err := nbus.Connect(cmd.Nats)
	if err != nil {
		return err
	}
	defer func() {
		if err := nb.Close(); err != nil {
			slog.Error("Failed to close NBus", "error", err)
		}
	}()

	col, err := collector.New(cmd.BasePath, a.Publish)
	if err != nil {
		return err
	}
	defer col.Close()

	stop := waitfor()
	<-stop
	slog.Info("Shutting down ged-journal")
	return nil
}
