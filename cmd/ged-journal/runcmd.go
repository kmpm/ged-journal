package main

import (
	"log/slog"

	"github.com/kmpm/ged-journal/public/fileapi"
)

type RunCmd struct {
	Nats        string `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string `help:"Nats context" default:""`
}

func (cmd *RunCmd) Run(cc *clicontext) error {
	slog.Info("Running file api")
	nc, err := connect(cmd.Nats, cmd.NatsContext)
	if err != nil {
		return err
	}
	a, err := fileapi.New(cc.BasePath, nc)
	if err != nil {
		return err
	}
	defer a.Close()

	slog.Info("Status", "status", a.Status)
	stop := waitfor()
	<-stop
	slog.Info("Shutting down ged-journal")
	return nil
}
