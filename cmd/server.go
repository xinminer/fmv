package cmd

import (
	"fmv/pkg/server"
	"github.com/urfave/cli/v2"
)

func init() {
	registerCommand(serverCmd)
}

var serverCmd = &cli.Command{
	Name:    "server",
	Aliases: []string{"srv"},
	Usage:   "start a server that receives files and listens on a specified port.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "specify a listening port",
			Value: "0.0.0.0:9988",
		},
		&cli.IntFlag{
			Name:  "chunk",
			Usage: "Size in MB of chunks size to be used as the streaming buffer (bigger might improve performance)",
			Value: 100,
		},
		&cli.StringFlag{
			Name:  "consul",
			Usage: "consul address",
		},
		&cli.StringFlag{
			Name:  "route",
			Usage: "route",
		},
	},
	Action: func(ctx *cli.Context) error {
		addr := ctx.String("addr")
		chunk := ctx.Int("chunk")
		consulAddr := ctx.String("consul")
		route := ctx.String("route")

		dests := ctx.Args().Slice()

		return server.StartServer(addr, chunk, dests, consulAddr, []string{route})
	},
}
