package main

import (
	"fmt"
	"runtime"
)

// go build cross.go
// GOARCH=386 go build cross.go
func main() {
	fmt.Println(runtime.GOOS, runtime.GOARCH)
}
