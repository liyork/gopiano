package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// mac 上ps一定要grep signalpart而不是signalpart.go
// 监听指定信号
// go run signalpart.go
// ctrl+c  ,interrupt
// kill pid  ,terminated
// kill -USR1 pid, user defined signal 1
// kill -USR2 pid, user defined signal 2
func main() {
	//合建chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c,kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
}
