package cmd

import (
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

func registerCommand(cmd *cli.Command) {
	commands = append(commands, cmd)
}

func NewApp() *cli.App {
	app := &cli.App{
		Name:                 "fmv",
		Usage:                "big file transfer",
		Commands:             commands,
		EnableBashCompletion: true,
	}

	return app
}
