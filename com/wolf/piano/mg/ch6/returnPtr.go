package main

import "fmt"

func returnPtr(x int) *int {
	y := x * x
	// y变量的内存地址
	return &y
}

// go run returnPtr.go
func main() {
	sq := returnPtr(10)
	// 解引用指针变量,返回存储在内存地址里的实际值
	fmt.Println("sq:", *sq)

	fmt.Println("sq:", sq)
}
