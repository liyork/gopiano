package main

import (
	"flag"
	"fmt"
	"time"
)

var period = flag.Duration("period", 1*time.Second, "sleep period")

// go run sleep.go
// go run sleep.go -period=2s
func main() {
	flag.Parse()
	fmt.Printf("Sleepingg for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}
