package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/kmpm/ged-journal/internal/misccli"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Cli struct {
	Loglevel  string `help:"Set log level" default:"info" short:"l" enum:"debug,info,warn,error"`
	Logfile   string `help:"Log to file" short:"f"`
	LogSource bool   `help:"Add source to log output"`

	Metrics   string       `help:"Enable prometheus metrics on address" short:"m" default:""`
	Collect   CollectCmd   `cmd:"" help:"Run the program"`
	Ls        LsCmd        `cmd:"" help:"List files in base-path"`
	Subscribe SubscribeCmd `cmd:"" aliases:"sub" help:"Subscribe to journal events"`
	Agent     AgentCmd     `cmd:"" help:"Run the agent"`
}

type clicontext struct {
	Username string
	HomePath string
	Metrics  bool
}

func main() {
	var cli Cli
	currUser, _ := user.Current()
	homeDir := currUser.HomeDir
	basePath := filepath.FromSlash(homeDir + "/Saved Games/Frontier Developments/Elite Dangerous")

	cc := clicontext{
		Username: currUser.Username,
		HomePath: homeDir,
	}

	ctx := kong.Parse(&cli, kong.Vars{"basepath": basePath})

	misccli.SetupLogging(cli.Loglevel, cli.Logfile, cli.LogSource)
	if cli.Metrics != "" {
		go func() {
			slog.Info("metrics enabled", "address", cli.Metrics)
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(cli.Metrics, nil)
			if err != nil {
				slog.Error("metrics server error", "error", err)
			}
		}()
		cc.Metrics = true
	} else {
		slog.Info("metrics disabled")
	}

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
