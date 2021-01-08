package bfbczd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWatigroup1(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1) //1
	go func() {
		defer wg.Done() //2
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1) //1
	go func() {
		defer wg.Done() //2
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	wg.Wait() //3
	fmt.Println("All goroutines complete.")
	//1.这里我们调用Add并传入参数1来表示一个goroutine正在开始。
	//2.在这里我们使用defer关键字来调用Done，以确保在退出goroutine的闭包之前，向WaitGroup表明了我们已经退出。
	//3.在这里，我们调用Wait，这将main goroutine，直到所有的goroutine都表明它们已经退出。
}

// 通常情况下，尽可能与要跟踪的goroutine就近且成对的调用Add，但有时候会一次性调用Add来跟踪一组goroutine
func TestWatigroup2(t *testing.T) {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}
