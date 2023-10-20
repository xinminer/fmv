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
			Name:  "dest",
			Usage: "upload dir or save path",
			Value: "/mnt/data1",
		},
	},
	Action: func(ctx *cli.Context) error {
		addr := ctx.String("addr")
		chunk := ctx.Int("chunk")
		dest := ctx.String("dest")
		server.StartServer(addr, chunk, dest)
		return nil
	},
}
