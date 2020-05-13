package base

import (
	"fmt"
	"testing"
	"unsafe"
)

// 宽度是数据类型的一个属性。Go程序中所有的实例都是一种数据类型，一个实例的宽度是由他的数据类型决定的，通常是8bit的整数倍。
// 通过unsafe.Sizeof()函数获取任何实例的宽度：

func TestBaseUnsafe(t *testing.T) {
	var s string
	var c complex128
	fmt.Println(unsafe.Sizeof(s))
	fmt.Println(unsafe.Sizeof(c))

	type S struct {
		a uint16
		b uint32
	}
	var s1 S
	fmt.Println(unsafe.Sizeof(s1))
}
