package main

import (
	"io"
	"log/slog"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/alecthomas/kong"
)

var globalLogLevel *slog.LevelVar

type cli struct {
	Loglevel  string `help:"Set log level" default:"info" short:"l" enum:"debug,info,warn,error"`
	Logfile   string `help:"Log to file" short:"f"`
	LogSource bool   `help:"Add source to log output"`

	Metrics   string       `help:"Enable prometheus metrics on address" short:"m" default:""`
	Collect   CollectCmd   `cmd:"" default:"1" help:"Run the program"`
	Ls        LsCmd        `cmd:"" help:"List files in base-path"`
	Subscribe SubscribeCmd `cmd:"" aliases:"sub" help:"Subscribe to journal events"`
}

type clicontext struct {
	Username string
	HomePath string
	Metrics  bool
}

// configure slog logging
func setupLogging(level, logfile string, source bool) {
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

	var w io.Writer
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("failed to open logfile", "file", logfile, "error", err)
		}
		w = io.MultiWriter(os.Stdout, file)
	} else {
		w = os.Stdout
	}
	handler := slog.NewJSONHandler(w, opts)
	logger := slog.New(handler)
	// buildInfo, _ := debug.ReadBuildInfo()
	child := logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			// slog.String("go_version", buildInfo.GoVersion),
		),
	)
	// log := slog.NewLogLogger(handler, slog.LevelError)
	slog.SetDefault(child)
}

func main() {
	var cli cli
	currUser, _ := user.Current()
	homeDir := currUser.HomeDir
	basePath := filepath.FromSlash(homeDir + "/Saved Games/Frontier Developments/Elite Dangerous")

	cc := clicontext{
		Username: currUser.Username,
		HomePath: homeDir,
	}

	ctx := kong.Parse(&cli, kong.Vars{"basepath": basePath})

	setupLogging(cli.Loglevel, cli.Logfile, cli.LogSource)

	slog.Info("Starting ged-journal", "user", cc.Username, "basepath", cc.BasePath, "loglevel", cli.Loglevel, "logfile", cli.Logfile)
	slog.Debug("cli", "cli", cli)

	slog.Info("Starting ged-journal", "user", cc.Username, "loglevel", cli.Loglevel, "logfile", cli.Logfile, "metrics", cli.Metrics)

	err := ctx.Run(&cc)
	ctx.FatalIfErrorf(err)
}

func waitfor() chan bool {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()
	return done
}
