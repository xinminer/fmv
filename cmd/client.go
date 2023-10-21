package cmd

import (
	"context"
	"fmt"
	"fmv/pkg/client"
	"fmv/pkg/consul"
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
		consulAddr := ctx.String("consul")
		tag := ctx.String("tag")
		chunk := ctx.Int("chunk")
		path := ctx.String("path")
		suffix := ctx.String("suffix")
		parallel := ctx.Int("parallel")

		ch := make(chan struct{}, parallel)

		gtimer.AddSingleton(ctx.Context, time.Second, func(ctx context.Context) {
			ch <- struct{}{}
			list, err := gfile.ScanDirFile(path, suffix, false)
			if err != nil {
				fmt.Printf("error (%s) in obtaining file list", err.Error())
				return
			}

			if len(list) == 0 {
				time.Sleep(5 * time.Second)
				return
			}

			go func(file string) {
				se, err := consul.Discovery("fmv-server", consulAddr, tag)
				if err != nil {
					fmt.Printf("discovery error: %s", err.Error())
				}
				se, err = consul.Discovery("fmv-server", consulAddr, "")
				if err != nil {
					fmt.Printf("discovery error: %s", err.Error())
					return
				}
				fc := client.NewFileClient(file, se, chunk)
				fc.SendFile()
				<-ch
			}(list[0])
		})

		select {}
	},
}
