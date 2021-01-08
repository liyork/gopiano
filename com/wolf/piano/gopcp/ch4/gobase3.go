package main

import (
	"fmt"
	"time"
)

func main() {
	name := "Eric"
	go func() {
		fmt.Printf("Hello, %s!\n", name)
	}()
	//name = "Harry"
	//time.Sleep(time.Millisecond)

	// 换个位置
	time.Sleep(time.Millisecond)
	name = "Harry"
}
