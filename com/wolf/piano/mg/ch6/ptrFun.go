package main

import "fmt"

// 指针参数
func getPtr(v *float64) float64 {
	return *v * *v
}

// go run ptrFun.go
func main() {
	x := 12.2
	fmt.Println(getPtr(&x))
	x = 12
	fmt.Println(getPtr(&x))
}
