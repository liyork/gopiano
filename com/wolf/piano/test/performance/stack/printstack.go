package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func myroutine(routineid string) {
	fmt.Println("Entry routine:", routineid)
	time.Sleep(1000 * time.Second)
	fmt.Println("Exit  routine:", routineid)
}

// go run printstack.go
// kill -SIGUSR1 <pid>
func main() {
	fmt.Println("Entry routine: main")

	// Register signal handler
	setupSignalHandler()

	// Launch goroutine 1
	go myroutine("myroutine1")

	// Launch goroutine 2
	go func(routineid string) {
		fmt.Println("Entry routine:", routineid)
		time.Sleep(1000 * time.Second)
		fmt.Println("Entry routine:", routineid)
	}("myroutine2")

	// wait main routine
	fmt.Scanln()
	fmt.Println("Exit  routine: main")
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for range c {
			dumpStacks()
		}
	}()
}

func dumpStacks() {
	buf := make([]byte, 1024)

	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}

	fmt.Printf("=== BEGIN goroutine stack dump ===\n")
	fmt.Printf("%s", buf)
	fmt.Printf("=== END   goroutine stack dump ===\n")
}
