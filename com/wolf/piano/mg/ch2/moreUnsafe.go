package main

import (
	"fmt"
	"unsafe"
)

// go run moreUnsafe.go
func main() {
	array := [...]int{0, 1, -2, 3, 4}
	pointer := &array[0]
	fmt.Print(*pointer, " ")
	// 数组第一个元素地址+数组元素大小
	memoryAddress := uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	for i := 0; i < len(array)-1; i++ {
		pointer = (*int)(unsafe.Pointer(memoryAddress))
		fmt.Print(*pointer, " ")
		memoryAddress = uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	}

	fmt.Println()
	pointer = (*int)(unsafe.Pointer(memoryAddress))
	// 尝试访问无效的数组元素，程序并不会出错而是会返回一个随机的数字。
	fmt.Print("One more: ", *pointer, " \n")
	memoryAddress = uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	fmt.Print("Two more: ", *pointer, " \n")
	fmt.Println()
}
