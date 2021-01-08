package test

import (
	"fmt"
	"strings"
	"testing"
)

func square(n int) int     { return n * n }
func negative(n int) int   { return -n }
func product(m, n int) int { return m * n }

func add1(r rune) rune { return r + 1 }

func TestBase(t *testing.T) {
	// function value
	f := square
	fmt.Println(f(3))

	f = negative
	fmt.Println(f(3))
	fmt.Printf("%T\n", f)

	fmt.Println(strings.Map(add1, "HAL-9000"))

	var rmdir []func()
	for _, d := range []int{1, 2, 3} {
		//需要临时赋值，不然由于for公用一个变量地址，导致下面函数中用的是一个值
		dir := d
		rmdir = append(rmdir, func() { fmt.Println("dir:", dir) })
	}
	for _, d := range rmdir {
		d()
	}
}
