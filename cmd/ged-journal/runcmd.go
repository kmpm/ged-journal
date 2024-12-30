package main

import (
	"log/slog"

	"github.com/kmpm/ged-journal/public/fileapi"
)

type RunCmd struct {
}

func (cmd *RunCmd) Run(cc *clicontext) error {
	slog.Info("Running file api")
	a, err := fileapi.New(cc.BasePath)
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
