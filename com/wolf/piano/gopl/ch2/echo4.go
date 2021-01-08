package main

import (
	"flag"
	"fmt"
	"strings"
)

// go run echo4.go -n=true a b c
// go run echo4.go -s=xx a b c
// go run echo4.go -s / a b c
// go run echo4.go -help

var n = flag.Bool("n", false, "omit trailing newline") //省略空行
var sep = flag.String("s", " ", "separator")

func main() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep))
	if !*n { // 不省略
		fmt.Println()
	}
}
