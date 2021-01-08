package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// 使用time.After处理chan中的超时问题，一同处理waitgroup的问题

// true表明t超时，false表明5s+w完成
func timeout(w *sync.WaitGroup, t time.Duration) bool {
	temp := make(chan int)
	go func() { // 等待5s，等待w进行done
		time.Sleep(5 * time.Second)
		defer close(temp)
		w.Wait()
	}()

	// 等待temp或超时t
	select {
	case <-temp:
		return false
	case <-time.After(t):
		return true
	}
}

// go run timeOut2.go 10000  --10s
// go run timeOut2.go 100
func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		fmt.Println("Need a time duration!")
		return
	}

	var w sync.WaitGroup
	w.Add(1)
	t, err := strconv.Atoi(arguments[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	duration := time.Duration(int32(t)) * time.Millisecond
	fmt.Printf("Timeout period is %s\n", duration)

	if timeout(&w, duration) {
		fmt.Println("Timed out!")
	} else {
		fmt.Println("OK!")
	}

	w.Done()
	if timeout(&w, duration) {
		fmt.Println("Timed out!")
	} else {
		fmt.Println("OK!")
	}
}
