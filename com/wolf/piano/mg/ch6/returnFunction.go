package main

import "fmt"

// closure

// 返回函数，伴有指向funReturnFun中局部变量i的状态
func funReturnFun() func() int {
	i := 0
	return func() int {
		i++
		return i * i
	}
}

// go run returnFunction.go
func main() {
	i := funReturnFun()
	j := funReturnFun()

	fmt.Println("1:", i())
	fmt.Println("2:", i())
	fmt.Println("j1:", j())
	fmt.Println("j2:", j())
	fmt.Println("3:", i())
}
