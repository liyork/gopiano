package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var CLOSEA int32 = 0
var DATA = make(map[int]bool)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// 产生随机数
func first(min, max int, out chan<- int) {
	for {
		if atomic.LoadInt32(&CLOSEA) == 1 {
			close(out)
			return
		}
		out <- random(min, max)
	}
}

// 从in读取，判重，若重复则CLOSEA=true，写入out中，当in被close则对out进行close
func second(out chan<- int, in <-chan int) {
	for x := range in {
		fmt.Print(x, " ")
		_, ok := DATA[x]
		if ok {
			atomic.StoreInt32(&CLOSEA, 1)
		} else {
			DATA[x] = true
			out <- x
		}
	}
	fmt.Println()
	close(out)
}

func third(in <-chan int) {
	var sum int
	sum = 0
	for x2 := range in {
		sum = sum + x2
	}
	fmt.Printf("The sum of the random numbers is %d\n", sum)
}

// go run -race pipeline.go 1 10
// go run pipeline.go 1 10
// go run pipeline.go 10 20
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Need two integer paramters!")
		os.Exit(1)
	}
	n1, _ := strconv.Atoi(os.Args[1])
	n2, _ := strconv.Atoi(os.Args[2])

	if n1 > n2 {
		fmt.Printf("%d should be smaller than %d\n", n1, n2)
		return
	}

	rand.Seed(time.Now().UnixNano())
	A := make(chan int)
	B := make(chan int)

	go first(n1, n2, A)
	go second(B, A)
	third(B)
}
