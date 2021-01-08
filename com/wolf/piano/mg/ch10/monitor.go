package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var readValue = make(chan int)
var writeValue = make(chan int)

func set(newValue int) {
	writeValue <- newValue
}

func readFromChan() int {
	return <-readValue
}

// 单线程处理，读取writeValue并放入value中，或者放入readValue中
func monitor() {
	var value int

	for {
		select {
		case newValue := <-writeValue:
			value = newValue
			fmt.Printf("%d", value)
		case readValue <- value:
		}
	}
}

// go run monitor.go 5
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give an integer!")
		return
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Going to create %d random numbers.\n", n)
	rand.Seed(time.Now().Unix())
	go monitor()

	var w sync.WaitGroup
	for r := 0; r < n; r++ {
		w.Add(1)
		go func() {
			defer w.Done()
			set(rand.Intn(10 * n))
		}()
	}
	w.Wait()
	fmt.Printf("\nLast value: %d\n", readFromChan())
}