package base

import (
	"fmt"
	"testing"
)

// 一直等到包含defer语句的函数执行完毕时，延迟函数（defer后的函数）才会被执行，
// 而不管函数是通过return的正常结束，还是由于panic导致的异常结束。
func TestDefer(t *testing.T) {

	i := return1()
	fmt.Println("i:", i)
}

func return1() int {
	defer println(1111)
	fmt.Println("222")
	return 1
}

// defer声明和执行顺序相反
func TestDeferSeq(t *testing.T) {

	defer println(1111)
	defer println(222)
	fmt.Println("3333")
}

type passValue struct {
	a int
}

func TestPassValue(t *testing.T) {
	value := passValue{1}
	// 对闭包的参数传递是值传递
	defer func(pv passValue) {
		fmt.Println("passValue1", pv)
	}(value)

	// 闭包内使用外部变量是引用传递
	defer func() {
		fmt.Println("passValue2", value)
	}()

	value.a = 2
	fmt.Println("value:", value)
}

// 得先定义defer才能最后被调用
func TestDeferUnder(t *testing.T) {

	fmt.Println("1111")
	testP()
	defer println(2222)
}

func testP() {
	panic(33333)
}
