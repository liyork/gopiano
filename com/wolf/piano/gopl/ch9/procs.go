package main

import (
	"fmt"
	"time"
)

// GOMAXPROCS=1 go run procs.go
// GOMAXPROCS=2 go run procs.go
func main() {
	for {
		go fmt.Print(0)
		fmt.Print(1)
		time.Sleep(time.Second)
	}
}
