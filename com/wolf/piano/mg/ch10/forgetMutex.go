package main

import (
	"fmt"
	"sync"
)

var m sync.Mutex

func function() {
	m.Lock()
	fmt.Println("Locked!")
}

// 忘记解锁 `sync.Mutex` 互斥体会产生崩溃
// go run forgetMutex.go
func main() {
	var w sync.WaitGroup

	w.Add(1)
	go func() {
		defer w.Done()
		function()
	}()

	w.Add(1)
	go func() {
		defer w.Done()
		function()
	}()

	w.Wait()
}
