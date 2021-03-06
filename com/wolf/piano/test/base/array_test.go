package base

import (
	"fmt"
	"testing"
)

// Go语言自带了数组类型，而切片类型是基于数组类型的抽象。
//定义一个数组时，需要指定数组长度和数组中元素的类型，[4]int定义了长度为4的数组，其中的元素类型为int。
// 一个数组的长度是固定的；长度是数组类型的一部分([4]int和[5]int就是两个不同的类型)。
//在没有显式初始化时，数组默认会将元素初始化为0。
//在内存中，[4]int表示为顺序排列的4个整数值
//Go语言中的数组是一个值。数组变量表示整个数组，而不是指向数组第一个元素的指针（就像C语言那样）。
// 意味着，将一个数组当作一个参数传递时，会完全拷贝数组中的内容（如果不想完全拷贝数组，可以传一个指向数组的指针）。
//可以把数组当成这样一种结构，它具有索引，有着固定的大小，可以用来存储同类型的元素。
func TestArrayBase(t *testing.T) {
	// 声明式带有固定长度
	var a [4]int
	// 0下标为1，其他默认0
	a[0] = 1

	for _, x := range a {
		fmt.Println("x:", x)
	}

}

func Test_array(t *testing.T) {
	// 定义一：var 数组名  [元素个数]数据类型
	var arr [8]int = [8]int{1, 2} // 部分初始化，默认0
	fmt.Println("arr", arr)

	// 定义二：使用自动类型推导创建数组,数组名：=[元素个数]数据类型
	bar := [10]int{1, 2}
	fmt.Println("bar", bar)

	// 定义三：使用三个点的自动类型推导，可以根据元素个数创建数组，数组长度可以随意变化
	// 数组名：=[...]数据类型
	car := [...]int{1, 2}
	fmt.Println("car", car)

	// 数组遍历
	for i, v := range arr { //i表示下标，v表示数组的值
		fmt.Println(i, v)
	}
}

func TestArrayLen(t *testing.T) {
	var arr [8]int = [8]int{1, 2}
	fmt.Println("arr", arr)
	fmt.Println("len:", len(arr))
	fmt.Println("cap:", cap(arr))
}

func TestArraySlice(t *testing.T) {
	var arr [8]int = [8]int{1, 2}
	fmt.Println("arr", arr)

	// slice
	ints := arr[0:4]
	fmt.Printf("%T\n", ints) // 切片

	ints = append(ints, 5)
	fmt.Println("ints:", ints)

	//var target []int， 拷贝回给target扩容，所以原先引用不能获取实际值
	var target = make([]int, 5)
	copy(target, ints)
	fmt.Println("target:", target)
}

// 参见passvalue_test.go
// go语言则是在编译期间就确定其大小，然后始终是作为值传递的
