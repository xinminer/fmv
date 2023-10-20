package cmd

import (
	"context"
	"fmt"
	"fmv/pkg/client"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/urfave/cli/v2"
	"time"
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
			Name:  "consul",
			Usage: "specify a server address",
			Value: "127.0.0.1:9988",
		},
		&cli.StringFlag{
			Name:  "tag",
			Usage: "specify a server address",
		},
		&cli.IntFlag{
			Name:  "chunk",
			Usage: "Size in MB of chunks size to be used as the streaming buffer (bigger might improve performance)",
			Value: 100,
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "upload dir or save path",
		},
		&cli.StringFlag{
			Name:  "suffix",
			Usage: "file suffix",
		},
		&cli.IntFlag{
			Name:  "parallel",
			Usage: "parallel transmission",
			Value: 5,
		},
	},
	Action: func(ctx *cli.Context) error {
		consul := ctx.String("consul")
		tag := ctx.String("tag")
		chunk := ctx.Int("chunk")
		path := ctx.String("path")
		suffix := ctx.String("suffix")
		parallel := ctx.Int("parallel")

		gtimer.AddSingleton(ctx.Context, time.Second, func(ctx context.Context) {
			list, err := gfile.ScanDirFile(path, suffix, false)
			if err != nil {
				fmt.Printf("error (%s) in obtaining file list", err.Error())
				return
			}

			if len(list) == 0 {
				time.Sleep(5 * time.Second)
				return
			}

			if len(list) > parallel {
				list = list[:parallel-1]
			}
			client.StartClient(list, chunk, consul, tag)
		})

		select {}
	},
}
