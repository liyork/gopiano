package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

// 内部s为局部变量，其值通过函数返回值返回，s本身为一指针，因此其指向的内存地址不会是栈而是堆
func StudentRegister(name string, age int) *Student {
	s := new(Student) //局部变量s逃逸到堆

	s.Name = name
	s.Age = age

	return s
}

// go run -gcflags=-m escape1.go
// ./escape.go:12:10: new(Student) escapes to heap
//./escape.go:22:29: new(Student) escapes to heap
//./escape.go:23:14: "register:" escapes to heap
func main() {
	register := StudentRegister("Jim", 18)
	fmt.Println("register:", register)
}
