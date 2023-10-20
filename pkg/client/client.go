package client

import (
	"fmt"
	"fmv/pkg/consul"
	"os"
)

func StartClient(files []string, chunkSize int, consulAddr string, tag string) {

	if len(files) == 0 {
		fmt.Println("Error: At least one filename must be provided")
		os.Exit(1)
	}

	done := make(chan bool)
	for _, file := range files {
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
			fc := NewFileClient(file, se, chunkSize)
			fc.SendFile()
			done <- true
		}(file)
	}
	for i := 0; i < len(files); i++ {
		<-done
	}
}
