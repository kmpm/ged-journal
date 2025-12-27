package main

import (
	"log/slog"

	"codeberg.org/kmpm/ged-common/pkg/nbus"

	"github.com/kmpm/ged-journal/public/agent"
)

type AgentCmd struct {
	Nats     nbus.Config `help:"Nats server address" embed:"" prefix:"nats." envprefix:"NATS_"`
	NoStatus bool        `help:"Do handle status messages" default:"false"`
}

func (cmd *AgentCmd) Run(cc *clicontext) error {
	slog.Info("Running Agent")
	nb, err := nbus.Connect(cmd.Nats)
	if err != nil {
		panic(err)
	}
	ag, err := agent.New(nb.Conn(), "ged.")
	if err != nil {
		panic(err)
	}
	if !cmd.NoStatus {
		for stat := range ag.Status() {
			slog.Info("Received status", "status", stat)
		}
	}

	<-waitfor()
	slog.Info("Shutting down agent")
	ag.Close()
	return nil
}
