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
