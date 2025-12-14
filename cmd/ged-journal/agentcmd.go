package main

import (
	"log/slog"

	"github.com/kmpm/ged-journal/public/agent"
)

type AgentCmd struct {
	Nats Nats `help:"Nats server address" embed:"" prefix:"nats."`

	NoStatus bool `help:"Do handle status messages" default:"false"`
}

func (cmd *AgentCmd) Run(cc *clicontext) error {
	slog.Info("Running Agent")
	nc, err := connect(&cmd.Nats)
	if err != nil {
		panic(err)
	}
	ag, err := agent.New(nc, "ged.")
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
