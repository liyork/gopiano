package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := make(chan int, 1)

	// 1s一个随机数发送
	ticker := time.NewTicker(time.Second)
	go func() {
		// 这个ticker如何被终止后感知到?
		for _ = range ticker.C {
			select {
			case intChan <- 1:
			case intChan <- 2:
			case intChan <- 3:
			}
		}
		//close(intChan)
		fmt.Println("End. [sender]")
	}()

	// 接收chan进行汇总
	var sum int
	for e := range intChan {
		fmt.Printf("Received: %v\n", e)
		sum += e
		if sum > 10 {
			fmt.Printf("Got: %v\n", sum)
			ticker.Stop()
			close(intChan)
			//for {
			//	if _, ok := <-intChan; !ok {
			//		break
			//	}
			//}
			break
		}
	}
	fmt.Println("End. [receiver]")
	time.Sleep(2 * time.Second)
}
