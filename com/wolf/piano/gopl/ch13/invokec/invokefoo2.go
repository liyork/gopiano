package main

// #cgo LDFLAGS: -L ./ -lfoo
// #include <stdio.h>
// #include <stdlib.h>
// #include "foo.h"
import "C"
import (
	"fmt"
)

// 注：linux环境下实验

// 创建静态共享库libfoo.a：
// gcc -c foo.c
// ar rv libfoo.a foo.o  // 不创建则下错误
// # command-line-arguments
///usr/bin/ld: cannot find -lfoo
//collect2: ld returned 1 exit status

// 可选：创建一个动态共享库
// gcc -c foo.c
// gcc -shared -Wl,-soname,libfoo.so -o libfoo.so  foo.o

// go run invokefoo.go
func main() {
	fmt.Println(C.count)
	C.foo()
}
