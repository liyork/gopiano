package base

import (
	"fmt"
	"testing"
)

// Go语言的类型转换和类型断言:
//类型转换在编译期完成，包括强制转换和隐式转换
//类型断言在运行时确定，包括安全类型断言和非安全类型断言

// 使用第二个参数，若不匹配不会导致panic
func TestAsserBase(t *testing.T) {
	//var i interface{} = "TT"
	var i interface{} = 77
	if value, ok := i.(int); ok {
		fmt.Printf("类型匹配整型：%d\n", value)
	} else if value, ok := i.(string); ok {
		fmt.Printf("类型匹配字符串:%s\n", value)
	}

	b, _ := i.([]int) //通过断言实现类型转换
	for _, v := range b {
		fmt.Print(v, " ")
	}
	fmt.Println()
	return
}

// nil不能assert成其他类型
func TestAssertNil(t *testing.T) {
	a := int32(1)
	var b interface{} = a
	i2, ok := b.(int32)
	fmt.Println("i is int", i2, ok)

	var c interface{} = nil
	i2, ok = c.(int32)
	fmt.Println("nil is int", i2, ok)
}

func TestErrorUpgrate(t *testing.T) {
	err1 := testErr1()
	fmt.Println("err1 is nil", err1 == nil)
}

func testErr1() error {
	err2 := testErr2()
	fmt.Println("err2 is nil", err2 == nil)
	return err2
}

func testErr2() *myError {
	return nil
}

type myError struct {
	code int
	msg  string
}

func (e *myError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("my error, code:%v, msg:%v", e.code, e.msg)
}

type TypeA interface{}
type TypeB interface {
	Foo()
}

// 接口之间在编译期间可以确定的情况下可以使用隐式类型转换，当然也可以用强制类型转换（不常用），所有情况下都可以使用类型断言。
func TestInterfaceConvert(t *testing.T) {
	// 编译时无法确定能不能转换情况下，用断言
	var a TypeA
	var b = a.(TypeB)
	fmt.Println(b)
	// 编译时，可以确定情况下，直接用转换
	var c TypeB
	var d = TypeA(c)
	fmt.Println(d)
}

// 普通类型向接口转换，可以使用隐式类型转换
func TestType2InterfaceConvert(t *testing.T) {
	var s = "abc"
	var a TypeA
	a = s
	fmt.Println(a)
}

// 接口向普通类型转换，只能使用类型断言，不能用转换因为编译器无法确定，可见断言是运行时，而转换时编译时
func TestInterface2TypeConvert(t *testing.T) {
	var s string
	var a TypeA
	s = a.(string)
	fmt.Println(s)
}

// 综上，优先使用断言，因为有ok可以判断，其次使用转换，基本上只有类型到接口的自动转换外其他都可以直接用断言
