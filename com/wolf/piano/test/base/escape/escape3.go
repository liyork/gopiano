package main

import "fmt"

// 很多函数参数为interface类型，比如fmt.Println(a ...interface{})，编译期间很难确定其参数的具体类型，也会产生逃逸

// go run -gcflags=-m escape3.go
func main() {
	s := "Escape"
	fmt.Println(s)
}
