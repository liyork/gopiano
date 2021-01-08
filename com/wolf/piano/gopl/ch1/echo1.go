package main

import (
	"fmt"
	"os"
)

// go build echo1.go
// ./echo1 a b c
func main() {
	// 命令本身
	fmt.Println("os.Args[0]:", os.Args[0])
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}
