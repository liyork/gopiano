package main

import "fmt"

type first struct{}

func (a first) F() {
	a.shared()
}

func (a first) shared() {
	fmt.Println("This is shared() from first!")
}

// 嵌入，继承
type second struct {
	first
}

// 重写
func (a second) shared() {
	fmt.Println("This is shared() from second!")
}

// go run goCoIn.go
func main() {
	first{}.F()
	second{}.shared()
	i := second{}
	j := i.first
	j.F()
}
