package main

import "C"
import (
	"fmt"
)

// 如果一个Go函数想要被C语言调用，必须要先导出。在函数实现的开始部分，添加一行以export开头的注释，加上函数的名字，这样C程序才知道怎么使
//export PrintMessage
func PrintMessage() {
	fmt.Println("A Go function!")
}

//export Multiply
func Multiply(a, b int) int {
	return a * b
}

// 用Go代码生成一个C共享库
// go build -o usedByC.o -buildmode=c-shared usedByC.go
//产生```usedByC.h```和```usedByC.o```两个文件ls
func main() {

}
