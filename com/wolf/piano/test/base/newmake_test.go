package base

import (
	"fmt"
	"testing"
)

func TestNewMakeDiff(t *testing.T) {
	// 分配并未初始化
	var p *[]int = new([]int) // allocates slice structure; *p == nil; rarely useful
	// Idiomatic，带有初始化
	var v []int = make([]int, 100) // the slice v now refers to a new array of 100 ints

	fmt.Println("p:", p)
	fmt.Println("v:", v)
}
