package server

import (
	"fmt"
	"log"
	"net"
)

func StartServer(addr string, chunkSize int, dest string) {
	// Start listening
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to bind to address: ", addr)
	}

	fmt.Printf("EasyTransfer running on %s with chunk size of %dMB, destination folder: \"%s\"\t...\n", addr, chunkSize, dest)

	// Accept concurrent connections
	for {
		conn, err := l.Accept()
		log.Println("Connection established...")
		if err != nil {
			log.Fatal(err)
		}
		fs := NewFileServer(conn, dest, chunkSize)
		go fs.HandleFile()
	}
}
