package main

import (
	"fmt"
)

// 全局执行顺序：
// 1. 先执行全局变量
// 2. 执行init()--go run initseqtest.go init1.go，依据顺序执行init
// 3. 执行main方法
var age = test()

func test() int {
	fmt.Println("test()")
	return 90
}

func init() {
	fmt.Println("init......")
}

// init 函数，通常在init函数中完成初始化工作
func main() {
	fmt.Println("main......")

	fmt.Println("Age is value", Age, "Name is value", Name)
}
