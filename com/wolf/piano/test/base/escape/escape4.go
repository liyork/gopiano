package main

import "fmt"

// 闭包引用对象逃逸

// 函数中原本属于局部变量的a和b由于闭包的引用，不得不将二者放到堆上，以致产生逃逸
func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

// go run -gcflags=-m escape4.go
func main() {
	// 获取一个闭包,闭包引用了函数的局部变量a和b
	f := Fibonacci()

	for i := 0; i < 10; i++ {
		fmt.Printf("Fibonacci: %d\n", f())
	}
}
