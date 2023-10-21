package main

import (
	"fmt"
	"os"
)

func main() {
	entries, _ := os.ReadDir("/Users/twosson")
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
