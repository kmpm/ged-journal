package main

type CollectorCmd struct {
	Folder string `arg:"" help:"Path to application log files" type:"existingdir"`
	Prefix string `help:"Subject prefix" default:"ged.collector."`
}

func (cmd *CollectorCmd) Run(cli *Cli) error {
	return simulator(cli, cmd.Folder, cmd.Prefix)
}
