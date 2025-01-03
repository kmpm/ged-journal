package main

type CollectCmd struct {
	Folder string `arg:"" help:"Path to application log files" type:"existingdir"`
	Prefix string `help:"Subject prefix" default:"ged.collector."`
}

func (cmd *CollectCmd) Run(cli *Cli) error {
	return simulator(cli, cmd.Folder, cmd.Prefix)
}
