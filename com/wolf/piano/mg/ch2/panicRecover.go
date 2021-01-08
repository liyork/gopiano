package main

import "fmt"

// 根据我们的意愿结束而没有panicking，因为defer中使用的匿名函数控制了局面！
func a() {
	fmt.Println("Inside a()")
	defer func() {
		if c := recover(); c != nil {
			fmt.Println("Recover inside a()!")
		}
	}()
	fmt.Println("About to call b()")
	b()
	fmt.Println("b() exited!")
	fmt.Println("Exiting a()")
}

func b() {
	fmt.Println("Inside b()")
	panic("Panic in b()!")
	fmt.Println("Exiting b()")
}

// go run panicRecover.go
func main() {
	a()
	fmt.Println("main() ended!")
}
