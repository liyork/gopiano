package main

// #include <stdio.h>
// #include <stdlib.h>
// #include "foo.h"
import "C"
import (
	"fmt"
)

// cgo提供了一种机制：它能够根据import”C”中引入的头文件，自动找到相应的源文件进行编译链接
// 这种机制的调用，需要用到go build命令
// 注：linux环境下实验
// go build -o invokefoo3  --直接编译当前目录build directory，然后执行./invokefoo3
func main() {
	fmt.Println(C.count)
	C.foo()
}
