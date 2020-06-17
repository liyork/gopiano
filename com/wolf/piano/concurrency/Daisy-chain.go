package main

import "fmt"

func f(left, right chan int) {
	left <- 1 + <-right
}

// 依次连接通过chan，最后从最右边写入，最左边取出
func main() {
	const n = 10000
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan int) { c <- 1 }(right)
	fmt.Println(<-leftmost)
}
