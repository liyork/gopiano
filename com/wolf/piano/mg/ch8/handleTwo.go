package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// `handleTwo.go`处理的信号是`SIGINFO`和`SIGINT`，在Golang中称为`syscall.SIGINFO`和`os.Interrupt`。*如果你查看`os`包文档，会发现在所有系统上只保证存在两个`siganal`，分别是`syscall.SIGKILL`和`syscall.SIGINT`，在`Go`中也定义为`os.Kill`和`os.Interrupt`。*

func handleSignal(signal os.Signal) {
	fmt.Println("handleSignal() Caught:", signal)
}

// go build handleTwo.go
// ./handleTwo
// ps ax | grep ./handleTwo | grep -v grep

//.^CCaught: interrupt --ctrl+c
//.Caught: interrupt  --kill -s INT 47988
//handleSignal() Caught: information request  --kill -s INFO 47988
// kill -s USR1 47988  --忽略
//.Killed:9  --kill -9 47988  --杀掉
func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINFO)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case os.Interrupt:
				fmt.Println("Caught:", sig)
			case syscall.SIGINFO:
				handleSignal(sig)
				return
			}
		}
	}()

	for {
		fmt.Printf(".")
		time.Sleep(20 * time.Second)
	}
}
