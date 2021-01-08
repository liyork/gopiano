package main

import "fmt"

// 可以重新声明已经在外层代码块中声明过的标识符，外层被屏蔽在内层
var v = "1, 2, 3"

// v被声明3次
func main() {
	v := []int{1, 2, 3}
	if v != nil {
		var v = 123
		fmt.Printf("%v\n", v)
	}
}
