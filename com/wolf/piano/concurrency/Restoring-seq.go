package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str  string
	wait chan bool
}

// 每个goroutine等待信号
func main() {
	c := make(chan Message)

	for i := 0; i < 2; i++ {
		fmt.Println("i:", i)
		go func(i int) {
			waitForIt := make(chan bool) // Shared between all messages.
			c <- Message{fmt.Sprintf("%s: %d", "msg", i), waitForIt}
			fmt.Println("go msg:", fmt.Sprintf("%s: %d", "msg", i))
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			<-waitForIt
		}(i) // 要是不传入就闭包了...
	}

	msg1 := <-c
	fmt.Println(msg1.str)

	msg2 := <-c
	fmt.Println(msg2.str)

	msg1.wait <- true
	msg2.wait <- true
}

func boring3(msg string, waitForIt chan bool) <-chan Message {
	c := make(chan Message)
	go func() {
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s: %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c
}
