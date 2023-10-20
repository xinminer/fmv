package client

import (
	"fmt"
	"os"
)

func StartClient(files []string, addr string, chunkSize int) {

	if len(files) == 0 {
		fmt.Println("Error: At least one filename must be provided")
		os.Exit(1)
	}

	done := make(chan bool)
	for _, file := range files {
		go func(file string) {
			fc := NewFileClient(file, addr, chunkSize)
			fc.SendFile()
			done <- true
		}(file)
	}
	for i := 0; i < len(files); i++ {
		<-done
	}
}
