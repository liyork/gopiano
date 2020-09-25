package base

import (
	"fmt"
	"testing"
)

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
