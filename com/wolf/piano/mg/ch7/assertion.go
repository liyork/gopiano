package main

import "fmt"

// go run assertion.go
func main() {
	// myInt变量具有动态类型int和值123
	var myInt interface{} = 123
	k, ok := myInt.(int)
	if ok {
		fmt.Println("Success:", k)
	}
	v, ok := myInt.(float64)
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("Failed without panicking!")
	}

	// 变量i的类型为int，值为123，存储在myInt中。因此由于int满足myInt接口，并且myInt接口不需要实现接口函数，所以myInt.(int)的值是一个int值。
	i := myInt.(int)
	fmt.Println("No cheking:", i)
	// panic
	j := myInt.(bool)
	fmt.Println(j)
}
