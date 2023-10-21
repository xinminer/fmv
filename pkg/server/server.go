package server

import (
	"fmt"
	"fmv/pkg/consul"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/shirou/gopsutil/disk"
	"log"
	"net"
)

func StartServer(addr string, chunkSize int, destinations []string, consulAddr string, tags []string) error {

	// Start listening
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Printf("EasyTransfer running on %s with chunk size of %dMB ...\n", addr, chunkSize)

	address := gstr.Split(addr, ":")
	dis := consul.DiscoveryConfig{
		ID:      guid.S(),
		Name:    "fmv-server",
		Tags:    tags,
		Port:    gconv.Int(address[1]),
		Address: address[0],
	}

	if err := consul.RegisterService(consulAddr, dis); err != nil {
		return err
	}

	ch := make(chan struct{}, 5)
	// Accept concurrent connections
	for {
		ch <- struct{}{}
		conn, err := l.Accept()
		log.Println("Connection established...")
		if err != nil {
			log.Printf("accept error: %s", err.Error())
			continue
		}
		dest, err := GetMaxCapPath(destinations)
		if err != nil {
			log.Printf("get max cap path: %s", err.Error())
			break
		}
		fs := NewFileServer(conn, dest, chunkSize)
		go func() {
			fs.HandleFile()
			<-ch
		}()
	}

	return nil
}

func GetMaxCapPath(paths []string) (string, error) {
	var maxCap uint64 = 0
	var rsp string
	for _, path := range paths {
		usage, err := disk.Usage(path)
		if err != nil {
			fmt.Printf("disk usage error: %s", path)
			continue
		}
		cap := usage.Free
		if cap < 53687091200 {
			continue
		}
		if cap > maxCap {
			maxCap = cap
			rsp = path
		}
	}

	if maxCap == 0 {
		return "", fmt.Errorf("no space available")
	}
	return rsp, nil
}
