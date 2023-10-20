package cmd

import (
	"fmv/pkg/client"
	"github.com/urfave/cli/v2"
)

func init() {
	registerCommand(clientCmd)
}

var clientCmd = &cli.Command{
	Name:    "client",
	Aliases: []string{"cli"},
	Usage:   "start an upload client.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "specify a server address",
			Value: "127.0.0.1:9988",
		},
		&cli.IntFlag{
			Name:  "chunk",
			Usage: "Size in MB of chunks size to be used as the streaming buffer (bigger might improve performance)",
			Value: 100,
		},
	},
	Action: func(ctx *cli.Context) error {
		addr := ctx.String("addr")
		chunk := ctx.Int("chunk")

		files := ctx.Args().Slice()

		client.StartClient(files, addr, chunk)
		return nil
	},
}
