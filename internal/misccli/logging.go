package misccli

import (
	"io"
	"log/slog"
	"os"
)

var globalLogLevel *slog.LevelVar

// SetupLogging configure slog logging
func SetupLogging(level, logfile string, addSource bool) {
	globalLogLevel = &slog.LevelVar{}
	opts := &slog.HandlerOptions{
		Level:     globalLogLevel,
		AddSource: addSource,
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
		// w = io.MultiWriter(os.Stdout, file)
		w = file
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
