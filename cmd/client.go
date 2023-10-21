package cmd

import (
	"fmt"
	"fmv/pkg/client"
	"fmv/pkg/consul"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/urfave/cli/v2"
	"os"
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

		for {
			ch <- struct{}{}

			var fileName string
			entries, _ := os.ReadDir(path)
			for _, entry := range entries {
				name := entry.Name()
				if gstr.HasSuffix(name, suffix) {
					fileName = name
				}
			}

			if fileName == "" {
				fmt.Println("not found file")
				time.Sleep(5 * time.Second)
				continue
			}

			fileName = fmt.Sprintf("%s/%s", path, fileName)

			if err := gfile.Move(fileName, fmt.Sprintf("%s.%s", fileName, "fmv")); err != nil {
				continue
			}
			fileName = fmt.Sprintf("%s.%s", fileName, "fmv")

			fmt.Println("send file : " + fileName)

			go func(fileName string) {
				se, err := consul.Discovery("fmv-server", consulAddr, tag)
				if err != nil {
					fmt.Printf("discovery error: %s", err.Error())
				}
				se, err = consul.Discovery("fmv-server", consulAddr, "")
				if err != nil {
					fmt.Printf("discovery error: %s", err.Error())
					return
				}
				fc := client.NewFileClient(fileName, se, chunk)
				fc.SendFile()
				<-ch
			}(fileName)
		}
	},
}
