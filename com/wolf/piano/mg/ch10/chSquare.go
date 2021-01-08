package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var times int

// 构造新chan放入cc，从c中读取,加总后再放入c
func f1(cc chan chan int, signal chan bool) {
	c := make(chan int)
	cc <- c
	defer close(c)

	sum := 0
	select {
	case x := <-c:
		for i := 0; i <= x; i++ {
			sum = sum + i
		}
		c <- sum
	case <-signal:
		return
	}
}

// go run chSquare.go 4
// go run chSquare.go 6
func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		fmt.Println("Need just one integer argument!")
		return
	}

	times, err := strconv.Atoi(arguments[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	cc := make(chan chan int)

	// 调用f1后，从cc中读取ch然后放入i，再从c中读取sum值,之后close掉signal的chan
	for i := 1; i < times+1; i++ {
		signal := make(chan bool)
		go f1(cc, signal)
		ch := <-cc
		ch <- i
		for sum := range ch {
			fmt.Print("Sum(", i, ")=", sum)
		}
		fmt.Println()
		time.Sleep(time.Second)
		close(signal)
	}
}
