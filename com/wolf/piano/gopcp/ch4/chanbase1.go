package main

import (
	"fmt"
	"time"
)

var exchangeChan = make(chan string, 3)

func main() {
	signalChan := make(chan struct{}, 1)
	waitFinishChan := make(chan struct{}, 2)
	go func() { // 接收
		<-signalChan
		fmt.Println("Receive a sync signal and wait a second... [receiver]")
		time.Sleep(time.Second)
		for {
			if elem, ok := <-exchangeChan; ok { // 取出所有非零值后，ok为false
				fmt.Println("Received:", elem, "[receiver]")
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		waitFinishChan <- struct{}{}
	}()

	go func() { // 发送
		for _, elem := range []string{"a", "b", "c", "d"} {
			exchangeChan <- elem
			fmt.Println("Sent:", elem, "[sender]")
			if elem == "c" { //发送3个后，通知
				signalChan <- struct{}{}
				fmt.Println("Sent a sync signal. [sender]")
			}
		}
		fmt.Println("Wait 2 seconds... [sender]")
		time.Sleep(time.Second * 2)
		close(exchangeChan)
		waitFinishChan <- struct{}{}
	}()

	<-waitFinishChan
	<-waitFinishChan
}
