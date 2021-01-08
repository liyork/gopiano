package main

import (
	"fmt"
	"time"
)

func A(a, b chan struct{}) {
	<-a
	fmt.Println("A()!")
	time.Sleep(time.Second)
	close(b)
}

func B(a, b chan struct{}) {
	<-a
	fmt.Println("B()!")
	close(b)
}

func C(a chan struct{}) {
	<-a
	fmt.Println("C()!")
}

// go run defineOrder.go
func main() {
	x := make(chan struct{})
	y := make(chan struct{})
	z := make(chan struct{})

	go C(z)    // 等待z有值
	go A(x, y) // 等待x有值后，关闭y
	go C(z)    // 等待z有值
	go B(y, z) // 等待y有值后，关闭z
	go C(z)    // 等待z有值

	close(x) // 启动x
	time.Sleep(3 * time.Second)
}
