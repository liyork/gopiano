package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handle(signal os.Signal) {
	fmt.Println("Received:", signal)
}

// go build handleAll.go
// ./handleAll
// ps ax | grep ./handleAll | grep -v grep

// .Ignoring: hangup  --kill -s HUP 38995
//Handling syscall.SIGUSR2!  --kill -s USR2 38995
//Ignoring: user defined signal 1  --kill -s USR1 38995
//Received: interrupt  --kill -s INT 38995
//^CReceived: interrupt  --ctrl+c
//Received: terminated  --kill -s TERM 38995
func main() {
	sigs := make(chan os.Signal, 1)
	// 所有输入信号都将被处理
	// 可以在同一程序中使用不同的通道和相同的信号多次调用`signal.Notify`。在这种情况下，每个相关通道都将收到一份它要处理的信号副本！
	signal.Notify(sigs)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case os.Interrupt:
				handle(sig)
			case syscall.SIGTERM:
				handle(sig)
				os.Exit(0)
			case syscall.SIGUSR2:
				fmt.Println("Handling syscall.SIGUSR2!")
			default:
				fmt.Println("Ignoring:", sig)
			}
		}
	}()

	for {
		fmt.Printf(".")
		time.Sleep(30 * time.Second)
	}
}
