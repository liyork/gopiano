package main

import "fmt"

func printNumbers() {
	for i := 0; i < 5; i++ {
		defer func(n int) {
			fmt.Printf("%d", n)
		}(i)
	}
}

func printNumbers2() {
	for i := 0; i < 5; i++ {
		defer func(n int) {
			fmt.Printf("%d", n)
		}(i * 2)
	}
}

func main() {
	//printNumbers()
	printNumbers2()
}
