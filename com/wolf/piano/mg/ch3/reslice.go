package main

import "fmt"

func main() {
	s1 := make([]int, 5)
	reSlice := s1[1:3]
	fmt.Println(s1)

	fmt.Println(reSlice)

	reSlice[0] = -100
	reSlice[1] = 123456

	// 通过re-slicing操作得到的切片，与原切片指向同一片内存地址
	fmt.Println(s1)
	fmt.Println(reSlice)
}
