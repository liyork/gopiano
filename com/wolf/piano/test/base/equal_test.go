package base

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDeepEqualArray(t *testing.T) {
	// 顺序不同则不等。
	var a = []string{"1", "2", "3"}
	var b = []string{"1", "3", "2"}
	if !reflect.DeepEqual(a, b) {
		fmt.Println("not equal")
	}
	fmt.Println("equal")
}

func TestDeepEqualMap(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"a": 1, "c": 3, "b": 2} // 换顺序
	fmt.Println("m1 == nil?", m1 == nil)
	fmt.Println("m2 != nil?", m2 != nil)
	fmt.Println("reflect.DeepEqual(m1,m2) = ", reflect.DeepEqual(m1, m2))
	fmt.Println()
	m3 := map[string]int{"a": 1, "b": 2, "c": 3, "d": 1}
	fmt.Println("reflect.DeepEqual(m1,m3) = ", reflect.DeepEqual(m1, m3))
}

func TestDeepEqualInterface(t *testing.T) {
	m1 := map[string]interface{}{"a": "1", "b": 2, "c": 3}
	m2 := map[string]interface{}{"a": 1, "c": "3", "b": 2}
	fmt.Println(`reflect.DeepEqual(m1["a"],m2["a"]`, reflect.DeepEqual(m1["a"], m2["a"]))
	fmt.Println(`reflect.DeepEqual(m1["b"],m2["b"]`, reflect.DeepEqual(m1["b"], m2["b"]))
}
