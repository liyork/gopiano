package main

import (
	"fmt"
	"github.com/spf13/pflag"
)

var uds string

func init() {
	pflag.StringVar(&uds, "uds", "", "HTTP server on unix domain socket")
}

// ./pflagdemo --uds=cccc
func main() {
	pflag.Parse()
	fmt.Println("ip1:", uds)
}
