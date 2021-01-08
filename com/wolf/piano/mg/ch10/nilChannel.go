package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 不断接收，当超时时设定c=nil不在监控c
func add(c chan int) {
	sum := 0
	// 每秒执行
	t := time.NewTimer(time.Second)
	for {
		select {
		case input := <-c: // 通道读取
			sum = sum + input
		case <-t.C:
			c = nil
			fmt.Println(sum)
		}
	}
}

// 不断放入
func send(c chan int) {
	for {
		c <- rand.Intn(10)
	}
}

// go run nilChannel.go
// go run nilChannel.go
func main() {
	c := make(chan int)
	go add(c)
	go send(c)
	time.Sleep(3 * time.Second)
}
