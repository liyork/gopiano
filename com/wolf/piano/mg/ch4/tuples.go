package main

import "fmt"

// 返回了包含三个整数元素的元组
func retThree(x int) (int, int, int) {
	return 2 * x, x * x, -x
}

// go run tuples.go
func main() {
	fmt.Println(retThree(10))
	// 元组赋值
	n1, n2, n3 := retThree(20)
	fmt.Println(n1, n2, n3)

	// 依靠这种元组操作，我们无需借助`temp`变量就可以实现两个数字的交换
	n1, n2 = n2, n1
	fmt.Println(n1, n2)

	// 三值交换
	x1, x2, x3 := n1*2, n1*n1, -n1
	fmt.Println(x1, x2, x3)
}
