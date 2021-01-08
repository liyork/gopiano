package main

import (
	"fmt"
	"os"
)

// go run echo2.go a b c
func main() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}
