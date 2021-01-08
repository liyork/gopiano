package main

import "fmt"

func d1() {
	for i := 3; i > 0; i-- {
		defer fmt.Print(i, " ")
	}
}

// 在循环结束后，i的值为0，因为是0值使循环终止。在循环结束后会执行被延迟的匿名函数。因为它没有参数，将值为0的i进行三次打印输出！
func d2() {
	for i := 3; i > 0; i-- {
		defer func() {
			fmt.Print(i, " ")
		}()
	}
	fmt.Println()
}

func d3() {
	for i := 3; i > 0; i-- {
		defer func(n int) {
			fmt.Print(n, " ")
		}(i)
	}
}

// go run defer.go
func main() {
	d1()
	d2()
	fmt.Println()
	d3()
	fmt.Println()
}
