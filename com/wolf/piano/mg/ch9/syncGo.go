package main

import (
	"flag"
	"fmt"
	"sync"
)

// go run syncGo.go -n 30
// sync.WaitGroup中的state1字段是一个保存计数器的元素，该计数器根据sync.Add()和sync.Done()调用而增加和减少。
func main() {
	n := flag.Int("n", 20, "Number of goroutines")
	flag.Parse()
	count := *n
	fmt.Printf("Going to create %d goroutines.\n", count)

	var waitGroup sync.WaitGroup

	fmt.Printf("%#v\n", waitGroup)
	for i := 0; i < count; i++ {
		// 每次调用sync.Add()都会增加sync.WaitGroup变量中的计数器。需要注意的是，在go语句之前调用sync.Add(1)非常重要，以防止出现任何竞争条件
		waitGroup.Add(1)
		go func(x int) {
			defer waitGroup.Done()
			fmt.Printf("%d ", x)
		}(i)
	}

	fmt.Printf("%#v\n", waitGroup)
	waitGroup.Wait()
	fmt.Println("\nExiting...")
}
