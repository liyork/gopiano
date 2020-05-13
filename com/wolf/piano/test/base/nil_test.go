package base

import (
	"fmt"
	"testing"
)

func TestNil(t *testing.T) {
	strings := make(map[string]map[string][]string)
	fmt.Println("strings:", strings)
	// make只能初始化第一层map
	i := strings["a"]
	fmt.Println("i:", i)
	// range对于nil可以很好处理
	for k, v := range i {
		fmt.Println(k, v)
	}
	i3 := len(i)
	fmt.Println("i3:", i3)

	mapmap := make(map[string]string)
	q := mapmap["a"]
	fmt.Println("q:", q) // 空串

	slic := []string{}
	//s := slic[0] // 报错
	//fmt.Println("s:", s)
	// range对于nil可以很好处理
	for _, value := range slic {
		fmt.Println("value:", value)
	}

	var a []int
	for _, value := range a {
		fmt.Println("value:", value)
	}

	fmt.Println("111")
}
