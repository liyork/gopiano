package base

import (
	"fmt"
	"testing"
)

// 可变参数

// 多个参数将会生成一个切片传入，函数内部按照切片来处理
//参数可以直接传递一个切片,使用`slice...`来表示切片,,注意:切片传入时不会生成新的切片，也就是说函数内部使用的切片与传入的切片共享相同的存储空间,如果函数内部修改了切片，可能会影响外部调用的函数。
func printfVariadic(str ...string) {
	for _, s := range str {
		fmt.Printf("%s\n", s)
	}
}

//可变参函数的特征:
//- 可变参数必须在函数参数列表的尾部，即最后一个（如放前面会引起编译时歧义）；
//- 可变参数在函数内部是作为切片来解析的；
//- 可变参数可以不填，不填时函数内部当成`nil`切片处理；
//- 可变参数必须是相同类型的（如果需要是不同类型的可以定义为interface{}类型）

// 传递的切片引用会复制，但是指向底层数组还是一个
func TestPassParam(t *testing.T) {
	str := []string{"1", "2"}
	passParam(str...)
	fmt.Printf("22 %p,%v\n", &str, str)
}

func passParam(str ...string) {
	fmt.Printf("11 %p,%v\n", &str, str)
	str[0] = "3"
}
