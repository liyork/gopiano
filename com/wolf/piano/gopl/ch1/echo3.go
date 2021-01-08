package main

import (
	"fmt"
	"os"
	"strings"
)

// go run echo3.go a b c
func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
