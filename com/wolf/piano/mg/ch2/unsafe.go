package main

import (
	"fmt"
	"unsafe"
)

// Unsafe code**是一种绕过go类型安全和内存安全检查的Go代码。大多数情况，unsafe code是和指针相关的
// 任何go指针都可以转化为```unsafe.Pointer```指针
// ```unsafe.Pointer```类型的指针可以覆盖掉go的系统类型。这毫无疑问很快，但是如果不小心或者不正确使用的话就会很危险，它给了开发者更多选择去掌控数据。*
func main() {
	var value int64 = 5
	var p1 = &value
	// 创造一个int32的p2指针去指向一个int64的value变量
	var p2 = (*int32)(unsafe.Pointer(p1))
	// 解引用一个指针
	fmt.Println("*p1: ", *p1)
	fmt.Println("*p2: ", *p2)

	*p1 = 5434123412312431212
	fmt.Println(value)
	// 使用32-bit的指针无法存一个64-bit的整数型
	fmt.Println("*p2: ", *p2)
	*p1 = 54341234
	fmt.Println(value)
	fmt.Println("*p2: ", *p2)
}
