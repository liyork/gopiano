package main

import (
	"fmt"
	"os"
	"reflect"
)

// 尽管t1和t2类型都基于int类型，本质上也与int类型相同，但Go语言将它们视为完全不同的类型。在Go编译器解析后的内部表示分别是main.t1和main.t2
type t1 int
type t2 int

type aa struct {
	X    int
	Y    float64
	Text string
}

// 比较两个变量的每个字段的接口是否一致
func (a1 aa) compareStruct(a2 aa) bool {
	r1 := reflect.ValueOf(&a1).Elem()
	r2 := reflect.ValueOf(&a2).Elem()

	for i := 0; i < r1.NumField(); i++ {
		if r1.Field(i).Interface() != r2.Field(i).Interface() {
			return false
		}
	}
	return true
}

func printMethods(i interface{}) {
	r := reflect.ValueOf(i)
	t := r.Type()
	fmt.Printf("Type to examine: %s\n", t)

	for j := 0; j < r.NumMethod(); j++ {
		m := r.Method(j).Type()
		fmt.Println(t.Method(j).Name, "-->", m)
	}
}

// go run advRefl.go
func main() {
	x1 := t1(100)
	x2 := t2(100)
	fmt.Printf("The type of x1 is %s\n", reflect.TypeOf(x1))
	fmt.Printf("The type of x2 is %s\n", reflect.TypeOf(x2))

	var p struct{}
	r := reflect.New(reflect.ValueOf(&p).Type()).Elem()
	fmt.Printf("The type of r is %s\n", reflect.TypeOf(r))

	a1 := aa{1, 2.1, "A1"}
	a2 := aa{1, -2, "A2"}

	if a1.compareStruct(a1) {
		fmt.Println("Equal!")
	}

	if !a1.compareStruct(a2) {
		fmt.Println("Not Equal!")
	}

	var f *os.File
	printMethods(f)
}
