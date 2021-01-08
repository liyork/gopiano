package main

import (
	"fmt"
	"os"
	"time"
)

// 提前终止
func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown. Press return to abort")
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("natural finish")
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	fmt.Println("lanunch finish")
}
