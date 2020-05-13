package main

import (
	"fmt"
	"testing"
)

type Delta struct {
	name string
}

func main() {
	var a = map[string][]Delta{}
	// 会有默认值
	delta := a["abc"]

	fmt.Println("delta:", delta)

	// 使用默认值进行添加
	newDeltas := append(a["abc"], Delta{"test"})
	fmt.Println("newDeltas:", newDeltas)

}

func TestSubSlice(t *testing.T) {
	s := "abcdefg"            // index从0开始
	s = string([]byte(s)[:3]) // endIndex不包含
	fmt.Println(s)

	s = "a我cd"
	s = string([]rune(s)[:3]) // rune方式，一个中文就只站一个数组下标,不考虑unicode字节问题
	fmt.Println(s)
}
