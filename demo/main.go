package main

import (
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
	"os"
)

func main() {
	entries, _ := os.ReadDir("/Users/twosson")
	for _, _ = range entries {
		fmt.Println(grand.N(0, 10))
	}
}
