package file

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"os"

	"github.com/google/uuid"
)

// PathExists 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GenFileName() string {
	u := uuid.New()
	return u.String()
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
