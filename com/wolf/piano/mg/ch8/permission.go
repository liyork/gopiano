package main

import (
	"fmt"
	"os"
)

// go run permission.go adobegc.log
// go run permission.go /dev/random
// ls -l /dev/random adobegc.log
func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Printf("usage:permissions filename\n")
		return
	}

	filename := arguments[1]
	info, _ := os.Stat(filename)
	mode := info.Mode()
	fmt.Println(filename, "mode is", mode.String()[1:10])
}
