package base

import (
	"fmt"
	"testing"
)

// 传递的是值的拷贝，不会影响原值
func passValue1(a int, b int) {
	a = b
}

func TestPassValue1(t *testing.T) {
	x := 5
	y := 10
	passValue1(x, y)
	fmt.Println(x)
}

// 传递地址的拷贝
func passPoint(a *int, b *int) {
	*a = *b
}

func TestPassPoint(t *testing.T) {
	x := 5
	y := 10
	passPoint(&x, &y)
	fmt.Println(x)
}

// go的数组传递也是值传递，进行拷贝
func passArray(a [1]int) {
	a[0] = 1111
}

func passArrayPointer(a *[1]int) {
	(*a)[0] = 1111
}

func TestPassArray(t *testing.T) {
	arr := [1]int{1}
	passArray(arr)
	fmt.Println(arr)
}

// 传递数组的引用
func TestPassArrayPointer(t *testing.T) {
	arr := [1]int{1}
	passArrayPointer(&arr)
	fmt.Println(arr)
}

func TestPassSlice(t *testing.T) {
	arr := []int{1}
	passSlice(arr)
	fmt.Printf("in TestPassSlice, p:%p, v:%v\n", &arr, arr)
}

// slice传递的是切片拷贝，但是底层用的一个数组
func passSlice(a []int) {
	// 修改底层array
	a[0] = 1111
	fmt.Printf("in passSlice, p:%p, v:%v\n", &a, a)
}
