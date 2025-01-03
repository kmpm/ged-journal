package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/alecthomas/kong"
)

var globalLogLevel *slog.LevelVar

type Cli struct {
	Loglevel string `help:"Set log level" default:"info" short:"l" enum:"debug,info,warn,error"`

	Delay       time.Duration `help:"Delay between events" default:"1s"`
	Nats        string        `help:"Nats server address" default:"nats://localhost:4222"`
	NatsContext string        `help:"Nats context" default:""`
	Collect     CollectCmd    `cmd:"" help:"Simulate ged-journal Collect"`
}

func setupLogging(level string, source bool) {
	globalLogLevel = &slog.LevelVar{}
	opts := &slog.HandlerOptions{
		Level:     globalLogLevel,
		AddSource: source,
	}
	switch level {
	case "debug":
		globalLogLevel.Set(slog.LevelDebug)
	case "info":
		globalLogLevel.Set(slog.LevelInfo)
	case "warn":
		globalLogLevel.Set(slog.LevelWarn)
	case "error":
		globalLogLevel.Set(slog.LevelError)
	default:
		globalLogLevel.Set(slog.LevelInfo)
		slog.Error("invalid log level", "level", level)
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
func main() {
	var cli Cli
	ctx := kong.Parse(&cli)
	setupLogging(cli.Loglevel, false)
	err := ctx.Run(cli)
	ctx.FatalIfErrorf(err)
}
