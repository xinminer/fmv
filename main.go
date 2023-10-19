package main

import (
	"fmv/cmd"
	"log"
	"os"
)

func main() {
	if err := cmd.NewApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
