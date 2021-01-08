package main

import (
	"fmt"
	"github.com/liyork/gopiano/com/wolf/piano/mg/ch6/aPackage"
)

// 需要创建一个名为main的package 有 `main()` 函数的包来创建执行文件
func main() {
	fmt.Println("Using aPackage!")
	aPackage.A()
	aPackage.B()
	fmt.Println(aPackage.MyConstant)
}
