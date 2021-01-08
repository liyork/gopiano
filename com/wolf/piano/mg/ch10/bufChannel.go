package main

import "fmt"

// go run bufChannel.go
func main() {
	numbers := make(chan int, 5)
	counter := 10

	// 只能放入5个，其余Not enough space for
	for i := 0; i < counter; i++ {
		select {
		case numbers <- i:
		default:
			fmt.Println("Not enough space for", i)
		}
	}

	// 只能读取5个,其余Nothing more to be don
	for i := 0; i < counter+5; i++ {
		select {
		case num := <-numbers:
			fmt.Println(num)
		default:
			fmt.Println("Nothing more to be done!")
			break
		}
	}
}
