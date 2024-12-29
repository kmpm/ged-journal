package main

import (
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"runtime/debug"

	"github.com/alecthomas/kong"
)

type cli struct {
	Loglevel string `help:"Set log level" default:"info" enum:"debug,info,warn,error"`
	Metrics  string `help:"Enable prometheus metrics on address" default:""`
	Run      RunCmd `cmd:"" default:"1" help:"Run the program"`
	Ls       LsCmd  `cmd:"" help:"List files in logpath"`
}

type clicontext struct {
	Username string
	HomePath string
	LogPath  string
}

// configure slog logging
func setupLogging(level string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	switch level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
		slog.Error("invalid log level", "level", level)
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	buildInfo, _ := debug.ReadBuildInfo()
	child := logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", buildInfo.GoVersion),
		),
	)
	// log := slog.NewLogLogger(handler, slog.LevelError)
	slog.SetDefault(child)
}

func main() {
	var cli cli
	currUser, _ := user.Current()
	homeDir := currUser.HomeDir
	cc := clicontext{
		Username: currUser.Username,
		HomePath: homeDir,
		LogPath:  filepath.FromSlash(homeDir + "/Saved Games/Frontier Developments/Elite Dangerous"),
	}

	ctx := kong.Parse(&cli, kong.Vars{"logpath": cc.LogPath})

	setupLogging(cli.Loglevel)

	err := ctx.Run(&cc)
	ctx.FatalIfErrorf(err)
}
