package init

import (
	"fmt"
	"testing"
)

// 先执行全局变量，再执行init(),再执行main方法
var age = test()

func test() int {
	fmt.Println("test()")
	return 90
}

func init() {
	fmt.Println("init......")
}

func TestInitSeq(t *testing.T) {
	fmt.Println("TestInitSeq......")
	fmt.Println("Age is value", Age, "Name is value", Name)
}

// 每个包可以包含任意多个 init 函数，这些函数都会在程序执行开始的时候被调用。所有被编译器发现的 init 函数都会安排在 main 函数之前执行
func TestInitBase(t *testing.T) {
	fmt.Println("444")
}

// 一个包中的多个 init 函数的执行是不确定的，不同包的 init 函数执行按照包导入的依赖关系决定执行顺序
func TestMultiInitSeq(t *testing.T) {
	fmt.Println("3333")
}
