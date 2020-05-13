package main

import (
	"github.com/astaxie/beego/logs"
	"time"
)

// todo 为什么没生成???，[2]    24293 killed     env GOTRACEBACK=crash ./testcoredump
// GOTRACEBACK变量可以控制程序在崩溃时，stack的输出情况
// ulimit -c unlimited
// go build testcoredump.go
// GOTRACEBACK=crash ./testcoredump
// 最后产生会coredump文件
func main() {

	logs.Info("Start...")
	defer logs.Info("exit.")
	i := 0
	c := make(chan int, 1)
	for {
		go func(i int) {
			mem := make([]int, 100*1024*1024)
			logs.Info("i=%d,mem:%p", i, mem)
			mem[0] = <-c
		}(i)
		i++
		time.Sleep(200 * time.Microsecond)
	}
}
