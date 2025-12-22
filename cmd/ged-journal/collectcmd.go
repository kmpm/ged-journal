package main

import (
	"log/slog"

	"github.com/kmpm/ged-journal/public/abus"
	"github.com/kmpm/ged-journal/public/collector"
)

type CollectCmd struct {
	BasePath string      `arg:"" help:"Path to application log files" default:"${basepath}"`
	Nats     abus.Config `embed:"" prefix:"nats." envprefix:"NATS_"`
}

func (cmd *CollectCmd) Run(cc *clicontext) error {
	slog.Info("Running Collect")
	a, err := abus.Connect(cmd.Nats)
	if err != nil {
		return err
	}
	defer func() {
		if err := a.Close(); err != nil {
			slog.Error("Failed to close ABus", "error", err)
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
