package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	v1    int
)

func change(i int) {
	mutex.Lock()
	time.Sleep(time.Second)
	v1 = v1 + 1
	if v1%10 == 0 {
		v1 = v1 - 10*i
	}
	mutex.Unlock()
}

func read() int {
	mutex.Lock()
	a := v1
	mutex.Unlock()
	return a
}

// go run mutex.go 2
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give me an integer!")
		return
	}

	numGR, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var waitGroup sync.WaitGroup

	fmt.Printf("%d ", read())
	for i := 0; i < numGR; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			change(i)
			fmt.Printf("-> %d", read())
		}(i)
	}

	waitGroup.Wait()
	fmt.Printf(", final -> %d\n", read())
}
