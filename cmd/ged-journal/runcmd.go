package main

import (
	"log/slog"

	"github.com/kmpm/ged-journal/internal/files"
)

type RunCmd struct {
}

func (cmd *RunCmd) Run(cc *clicontext) error {
	slog.Info("Starting ged-journal", "user", cc.Username, "logpath", cc.LogPath)
	m := files.New(cc.LogPath)
	stat, err := m.GetStatus()
	if err != nil {
		return err
	}
	slog.Info("Status", "status", stat)
	return nil
}
