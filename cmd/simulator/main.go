package main

import (
	"log/slog"
	"time"

	"github.com/alecthomas/kong"
	"github.com/kmpm/ged-journal/internal/misccli"
)

var globalLogLevel *slog.LevelVar

type Cli struct {
	Loglevel    string        `help:"Set log level" default:"info" short:"l" enum:"debug,info,warn,error"`
	Logfile     string        `help:"Log to file" short:"f"`
	Delay       time.Duration `help:"Delay between events" default:"1s"`
	Nats        string        `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string        `help:"Nats context" default:""`
	NatsCreds   string        `help:"Nats credentials" default:"" type:"existingfile"`
	Collect     CollectCmd    `cmd:"" help:"Simulate ged-journal Collect"`
}

func main() {
	var cli Cli
	ctx := kong.Parse(&cli)
	misccli.SetupLogging(cli.Loglevel, cli.Logfile, false)
	err := ctx.Run(cli)
	ctx.FatalIfErrorf(err)
}
